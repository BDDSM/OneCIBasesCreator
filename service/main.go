// +build windows

package main

import (
	"OneCIBasesCreator/args"
	"fmt"
	"golang.org/x/sys/windows/svc"
	"log"
	"os"
	"strings"
)

func main() {

	var err error

	if len(os.Args) < 2 {
		handleError(errBadCommand("no command specified"))
	}

	instance, err := args.Instance()
	handleError(err)

	isIntSess, err := svc.IsAnInteractiveSession()
	if err != nil {
		log.Fatalf("failed to determine if we are running in an interactive session: %v", err)
	}
	if !isIntSess {
		runService(instance, false)
		return
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "debug":
		runService(instance, true)
		return
	case "install":
		err = installService(instance, os.Args[2:])
	case "remove":
		err = removeService(instance)
	case "start":
		err = startService(instance)
	case "stop":
		err = controlService(instance, svc.Stop, svc.Stopped)
	case "pause":
		err = controlService(instance, svc.Pause, svc.Paused)
	case "continue":
		err = controlService(instance, svc.Continue, svc.Running)
	default:
		handleError(errBadCommand(fmt.Sprintf("invalid command %s", cmd)))
	}
	handleError(err)
	return
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		args.Usage()
	}
}

func errBadCommand(errmsg string) error {
	return fmt.Errorf(
		"%s\n\n"+
			"usage: %s <command>\n"+
			"       where <command> is one of\n"+
			"       install, remove, debug, start, stop, pause or continue.\n",
		errmsg, os.Args[0])
}
