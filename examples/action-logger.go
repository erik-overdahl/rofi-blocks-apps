package examples

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/erik-overdahl/rofi-blocks-apps/pkg/rofi"
	"github.com/erik-overdahl/rofi-blocks-apps/pkg/apps"
)

type ActionLoggerApp struct {
	apps.AppBase
	lineNum int
}

func MakeActionLoggerApp() *ActionLoggerApp {
	app := &ActionLoggerApp{
		lineNum: 0,
	}
	app.AppBase = apps.MakeApp(app.handleEvent, app.loop)
	return app
}

func (app *ActionLoggerApp) Name() string {
	return "action-logger"
}

func (app *ActionLoggerApp) Start() error {
	initial := []rofi.OutputUpdate{
		rofi.PromptUpdate{"Updating input also logs action"},
		rofi.MessageUpdate{fmt.Sprintf("Time: %s", time.Now().Format("03:04:05"))},
		rofi.InputActionUpdate{rofi.INPUT_ACTION_SEND},
	}
	app.SendOutput(initial)
	return app.AppBase.Start()
}

func (app *ActionLoggerApp) ShouldReceiveInBackground() bool {
	return true
}

func (app *ActionLoggerApp) handleEvent(event rofi.RofiBlocksEvent) error {
	text, err := json.Marshal(event)
	if err != nil {
		return err
	}
	eventLine := rofi.NewRofiBlocksLine()
	eventLine.Text = string(text)
	eventLine.Data = fmt.Sprintf("%d", app.lineNum)
	app.SendOutput([]rofi.OutputUpdate{
		rofi.AddLineUpdate{Prepend: true, Line: *eventLine},
	})
	app.lineNum++
	return nil
}

func (app *ActionLoggerApp) loop() {
	for {
		select {
		case <-app.Done():
			return
		case <-time.After(time.Second):
			app.SendOutput([]rofi.OutputUpdate{
				rofi.MessageUpdate{fmt.Sprintf("Time: %s", time.Now().Format("03:04:05"))},
			})
		}
	}
}
