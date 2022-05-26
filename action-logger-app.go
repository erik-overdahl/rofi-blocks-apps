package main

import (
	"fmt"
	"time"
)

type ActionLoggerApp struct {
	updates chan []OutputUpdate
	lineNum int
}

func (app *ActionLoggerApp) Name() string {
	return "action logger"
}

func (app *ActionLoggerApp) Init(updates chan []OutputUpdate) {
	initial := []OutputUpdate{
		PromptUpdate{value: "Updating input also logs action"},
		MessageUpdate{value: fmt.Sprintf("Time: %s", time.Now().Format("03:04:05"))},
		InputActionUpdate{value: INPUT_ACTION_SEND},
	}
	app.updates = updates
	app.updates <- initial
}

func (app *ActionLoggerApp) HandleEvent(event RofiBlocksEvent) *ProcessChange {
	newLine := &RofiBlocksLine{
		Text: fmt.Sprintf(`{"name":"%s","data":"%s","value":"%s"}`, event.Name, event.Data, event.Value),
		Data: fmt.Sprintf("%d", app.lineNum),
	}
	app.updates <- []OutputUpdate{
		AddLineUpdate{prepend: true, line: newLine},
	}
	app.lineNum++
	return nil
}

func (app *ActionLoggerApp) Loop() {
	for {
		time.Sleep(1 * time.Second)
		app.updates <- []OutputUpdate{
			MessageUpdate{value: fmt.Sprintf("Time: %s", time.Now().Format("03:04:05"))},
		}
	}
}
