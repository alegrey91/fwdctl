package daemon

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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
func Start() {
	infoLogger.Println(banner())

	err := createPidFile()
	if err != nil {
		errorLogger.Println(err)
		os.Exit(1)
	}
	defer removePidFile()

	sigChnl := make(chan os.Signal, 1)
	signal.Notify(sigChnl, syscall.SIGTERM)
	exitcChnl := make(chan bool, 1)

	go func() {
		for {
			time.Sleep(time.Second)
			select {
			case <-sigChnl:
				infoLogger.Println("daemon stopped")
				exitcChnl <- true
			default:
				infoLogger.Println("DO SOMETHING")
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
	syscall.Kill(pid, syscall.SIGTERM)
}
