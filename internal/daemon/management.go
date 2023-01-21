package daemon

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
func Start(rulesFile string) {
	// Current implementation sucks.
	// Code should be improved by removing
	// useless and redundant code.
	// Also the ruleSet.Diff() method,
	// should return a "path" to pass from
	// A to B, without applying any changes.
	// These should be done here.
	// Anyway, for now, it works.
	infoLogger.Println(banner())

	err := createPidFile()
	if err != nil {
		errorLogger.Println(err)
		os.Exit(1)
	}
	defer func() {
		err = removePidFile()
	}()
	infoLogger.Println("PID file created")

	// preparing rule set from rules file
	ruleSet, err := rules.NewRuleSetFromFile(rulesFile)
	if err != nil {
		errorLogger.Println(err)
		os.Exit(1)
	}
	// apply all the rules present in rulesFile
	for ruleId, rule := range ruleSet.Rules {
		err = iptables.CreateForward(rule.Iface, rule.Proto, rule.Dport, rule.Saddr, rule.Sport)
		if err != nil {
			infoLogger.Printf("rule %s - %v\n", ruleId, err)
		}
	}
	infoLogger.Println("rules from file have been applied")

	// preparing viper module to manage rules file
	v := viper.New()
	v.SetConfigFile(rulesFile)
	v.OnConfigChange(func(e fsnotify.Event) {
		infoLogger.Println("configuration has changed")
		newRuleSet, err := rules.NewRuleSetFromFile(rulesFile)
		if err != nil {
			errorLogger.Println(err)
			os.Exit(1)
		}
		err = newRuleSet.Diff(ruleSet)
		if err != nil {
			errorLogger.Println(err)
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
			time.Sleep(time.Second)
			select {
			case <-sigChnl:
				// flush rules before exit
				for _, rule := range ruleSet.Rules {
					err := iptables.DeleteForwardByRule(rule.Iface, rule.Proto, rule.Dport, rule.Saddr, rule.Sport)
					if err != nil {
						errorLogger.Println(err)
					}
				}
				infoLogger.Println("daemon stopped")
				exitcChnl <- true
			default:
			}
		}
	}()

	<-exitcChnl
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
