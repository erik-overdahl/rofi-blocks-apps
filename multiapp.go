package main

import (
	"strconv"
)

func getApp(appName string) App {
	switch appName {
	case "action-logger":
		return &ActionLoggerApp{}
	case "focus-lines":
		return &FocusLinesApp{}
	case "theme-selector":
		return &ThemeSelectorApp{}
	case "show-lines":
		return &ShowLinesApp{}
	case "volume":
		return &VolumeApp{}
	case "powermenu":
		return &PowermenuApp{}
	default:
		return &AppNotFound{command: appName}
	}
}

type appWrapper struct {
	app       App
	updates   chan []OutputUpdate
	lastState *RofiBlocksOutput
}

type MultiApp struct {
	updates chan []OutputUpdate
	swap    chan appWrapper
	apps    []appWrapper
	current appWrapper
}

func MakeMultiApp(apps ...string) App {
	children := make([]appWrapper, len(apps), len(apps))
	for i, appName := range apps {
		app := getApp(appName)
		appUpdates := make(chan []OutputUpdate)
		go app.Init(appUpdates)
		switch a := app.(type) {
		case LoopingApp:
			go a.Loop()
		}
		children[i] = appWrapper{
			app:       app,
			updates:   appUpdates,
			lastState: &RofiBlocksOutput{},
		}
	}
	return &MultiApp{apps: children}
}

func (app *MultiApp) Init(updates chan []OutputUpdate) {
	app.updates = updates
	swap := make(chan appWrapper)
	app.swap = swap
	app.current = app.apps[0]
}

func (app *MultiApp) HandleEvent(event RofiBlocksEvent) *ProcessChange {
	if event.Name == CUSTOM_KEY {
		appIdx, _ := strconv.Atoi(event.Value)
		if appIdx < len(app.apps) && appIdx > 0 {
			new := app.apps[appIdx-1]
			app.swap <- new
		}
	}
	switch a := app.current.app.(type) {
	case EventHandlingApp:
		return a.HandleEvent(event)
	default:
		return nil
	}
}

func (app *MultiApp) Loop() {
	for {
		select {
		case new := <-app.swap:
			app.updates <- []OutputUpdate{
				SnapshotState{app.current.lastState},
				RestoreState{new.lastState},
			}
			app.current = new
		case output := <-app.current.updates:
			app.updates <- output
		}
	}
}
