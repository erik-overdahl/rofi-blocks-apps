package main

import (
	"log"
	"os"

	"github.com/erik-overdahl/rofi-blocks-apps/pkg/rofi"
	"github.com/erik-overdahl/rofi-blocks-apps/pkg/apps"
	"github.com/erik-overdahl/rofi-blocks-apps/examples"
)

func main() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	rofiProcess, err := rofi.MakeRofiProcess()
	if err != nil {
		log.Fatalf("Unable to start rofi process: %v", err)
	}
	defer rofiProcess.Stop()

	app := apps.MakeMultiApp()
	app.Add(
		examples.MakeFocusLinesApp(),
		examples.MakeShowLinesApp(),
		examples.MakeActionLoggerApp(),
	)

	go rofiProcess.SendUpdates(app.Output())
	go rofiProcess.ReadEvents(app.Input())

	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start app: %#v", err)
	}

	if err := rofiProcess.Start(); err != nil {
		log.Fatalf("Rofi process failed to start: %v", err)
	}
	HandleRofiExit(rofiProcess.ListenProcessExit())
}

func HandleRofiExit(state *os.ProcessState, err error) {
	switch state.ExitCode() {
	case 0:
		if err == nil {
			os.Exit(0)
		}
	case 65:
		log.Println("Rofi displayed an error and user exited")
		os.Exit(0)
	default:
		os.Exit(1)
	}
}
