package daemon

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alegrey91/fwdctl/internal/hooks"
	"github.com/alegrey91/fwdctl/internal/rules"
	"github.com/alegrey91/fwdctl/pkg/iptables"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

func init() {
	// Initialize different loggers
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Print daemon banner
func banner() string {
	return `
┌─┐┬ ┬┌┬┐┌─┐┌┬┐┬  
├┤ │││ │││   │ │  
└  └┴┘─┴┘└─┘ ┴ ┴─┘
┌┬┐┌─┐┌─┐┌┬┐┌─┐┌┐┌
 ││├─┤├┤ ││││ ││││
─┴┘┴ ┴└─┘┴ ┴└─┘┘└┘`
}

// Start run fwdctl in daemon mode
// The flow has the following steps:
// - Open the rules file
// - Create all the defined forwards (running their on-create hooks)
// - Start listening on rules file changes
// - When file changes:
//   - Calculate the Diff between old and new ruleset
//   - Delete unwanted forwards (running their on-delete hooks)
//   - Create wanted forwards (running their on-create hooks)
//
// - Listen for SIGTERM signals to gracefuly shutdown
// - When SIGTERM signal occurs:
//   - Delete all the applied forwards (running their on-delete hooks)
//   - Shutdown the daemon
//
// When a hook fails and continueOnHookFail is false, the daemon logs the
// error, tears down its forwards and exits non-zero. When it is true, the
// failure is only logged and the daemon keeps running.
func Start(ipt *iptables.IPTablesInstance, rulesFile string, continueOnHookFail bool) int {
	infoLogger.Println(banner())

	err := createPidFile()
	if err != nil {
		errorLogger.Println(err)
		return 1
	}
	defer func() {
		err = removePidFile()
	}()
	infoLogger.Println("PID file created")

	// preparing rule set from rules file
	rulesContent, err := os.Open(rulesFile)
	if err != nil {
		errorLogger.Printf("error opening file: %v", err)
		return 1
	}
	ruleSet, err := rules.NewRuleSetFromFile(rulesContent)
	if err != nil {
		errorLogger.Println(err)
		return 1
	}

	// mu guards ruleSet, shared between the config-change watcher and the
	// SIGTERM handler.
	var mu sync.Mutex

	// hookFailed reports a hook error according to the failure policy.
	// It returns true when the daemon must shut down (policy is to fail
	// and a hook did fail).
	hookFailed := func(err error) bool {
		if err == nil {
			return false
		}
		if continueOnHookFail {
			infoLogger.Printf("ignoring hook failure: %v", err)
			return false
		}
		errorLogger.Printf("hook failed: %v", err)
		return true
	}

	// shutdown removes every applied forward, running each on-delete hook.
	shutdown := func() {
		mu.Lock()
		defer mu.Unlock()
		for _, rule := range ruleSet.Rules {
			rule := rule
			if err := ipt.DeleteForwardByRule(&rule); err != nil {
				errorLogger.Println(err)
				continue
			}
			hookFailed(hooks.RunOnDelete(&rule))
		}
	}

	// apply all the rules present in rulesFile
	for ruleId, rule := range ruleSet.Rules {
		rule := rule
		err = ipt.CreateForward(&rule)
		if err != nil {
			infoLogger.Printf("rule %s - %v\n", ruleId, err)
			continue
		}
		if hookFailed(hooks.RunOnCreate(&rule)) {
			shutdown()
			return 1
		}
	}
	infoLogger.Println("rules from file have been applied")

	// preparing viper module to manage rules file
	v := viper.New()
	v.SetConfigFile(rulesFile)
	v.OnConfigChange(func(e fsnotify.Event) {
		infoLogger.Println("configuration has changed")
		rulesContent, err := os.Open(rulesFile)
		if err != nil {
			errorLogger.Printf("error opening file: %v", err)
			return
		}
		newRuleSet, err := rules.NewRuleSetFromFile(rulesContent)
		if err != nil {
			errorLogger.Println(err)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		rsd := rules.Diff(ruleSet, newRuleSet)
		// delete all the rules to be removed
		for _, rule := range rsd.ToRemove {
			err = ipt.DeleteForwardByRule(rule)
			if err != nil {
				errorLogger.Println(err)
				continue
			}
			hookFailed(hooks.RunOnDelete(rule))
		}
		// create all the rules to be added
		for _, rule := range rsd.ToAdd {
			err = ipt.CreateForward(rule)
			if err != nil {
				errorLogger.Println(err)
				continue
			}
			hookFailed(hooks.RunOnCreate(rule))
		}
		// set the new rule set as the current one
		ruleSet = newRuleSet
	})
	v.WatchConfig()

	sigChnl := make(chan os.Signal, 1)
	signal.Notify(sigChnl, syscall.SIGTERM)
	exitcChnl := make(chan bool, 1)

	go func() {
		for {
			select {
			case <-sigChnl:
				// flush rules before exit, running on-delete hooks
				shutdown()
				infoLogger.Println("daemon stopped")
				exitcChnl <- true
			default:
				continue
			}
		}
	}()
	<-exitcChnl
	return 0
}

// Stop send a SIGTERM signal to the daemon process
func Stop() {
	infoLogger.Println("stopping daemon")
	pid, err := readPidFile()
	if err != nil {
		errorLogger.Println(err)
	}
	err = syscall.Kill(pid, syscall.SIGTERM)
	if err != nil {
		errorLogger.Println(err)
	}
}
