package main

import (
	"log"
	"os"
	"sync"
)

const ICONS_DIR = "/usr/share/icons/Adwaita/scalable"

var (
	logger log.Logger      = *log.New(os.Stderr, "", 4)
	wg     *sync.WaitGroup = &sync.WaitGroup{}
)

type ProcessChange struct {
	Args []string
}

func main() {
	events := make(chan RofiBlocksEvent)
	outputUpdates := make(chan []OutputUpdate)
	outputReceiver := make(chan RofiBlocksOutput)
	receiverReady := make(chan int)
	go accumulateUpdates(outputUpdates, outputReceiver, receiverReady)
	// read args to get which apps we are running
	// args := os.Args[1:]
	// INITIALIZE the app - should give us the initial output
	//   to send to rofi and /optionally/ arguments to the rofi process
	var app App
	if len(os.Args) > 2 {
		app = MakeMultiApp(os.Args[1:]...)
	} else {
		app = getApp(os.Args[1])
	}
	// app = &ThemeSelectorApp{}
	app.Init(outputUpdates)
	// start the rofi process
	// read Events from the process
	// call HandleEvent on app, passing in the Event
	// read OutputUpdates from a channel, apply them to the Output
	//   and write the Output to stdout
	rofi, err := StartRofiProcess()
	if err != nil {
		// TODO
		panic(err)
	}
	go sendOutput(rofi, outputReceiver, receiverReady)
	go readEvents(rofi, events)
	switch a := app.(type) {
	case LoopingApp:
		logger.Println("Starting loop")
		wg.Add(1)
		go func(app LoopingApp) {
			defer wg.Done()
			app.Loop()
			logger.Println("loop done")
		}(a)
	}
	switch a := app.(type) {
	case EventHandlingApp:
		for event := range events {
			logger.Printf("Received event: %+v\n", event)
			changeProcess := a.HandleEvent(event)
			if changeProcess != nil {
				rofi.Stop()
				rofi, err = StartRofiProcess(changeProcess.Args...)
				go sendOutput(rofi, outputReceiver, receiverReady)
				go readEvents(rofi, events)
			}
		}
	}
	wg.Wait()
	logger.Println("Exiting")
}

// gather OutputUpdates until the goroutine that applies them is ready
func accumulateUpdates(sender chan []OutputUpdate, receiver chan RofiBlocksOutput, receiverReady chan int) {
	output := RofiBlocksOutput{}
	isChanged := false
	ready := false
	receivingProcess := -1
	lastProcess := -1
	for {
		select {
		case r := <-receiverReady:
			// logger.Println("Got receiver ready")
			ready = true
			receivingProcess = r
		case updates := <-sender:
			// logger.Println("Got updates to apply to output")
			for _, update := range updates {
				update.Apply(&output)
			}
			isChanged = true
		default:
			if ready {
				if receivingProcess != lastProcess {
					output.ChangeActiveEntry = true
					output.ChangeInput = true
					output.ChangeInputAction = true
					output.ChangeLines = true
					output.ChangeMessage = true
					output.ChangeOverlay = true
					output.ChangePrompt = true
				}
				if (receivingProcess != lastProcess) || isChanged {
					logger.Printf("Sending to %d\n", receivingProcess)
					receiver <- output
					isChanged = false
					ready = false
					lastProcess = receivingProcess
					output.ResetFlags()
				}
			}

		}
	}
}
