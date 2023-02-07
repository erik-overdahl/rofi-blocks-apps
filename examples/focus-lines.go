package examples

import (
	"strconv"

	"github.com/erik-overdahl/rofi-blocks-apps/pkg/rofi"
	"github.com/erik-overdahl/rofi-blocks-apps/pkg/apps"
)

type FocusLinesApp struct {
	apps.AppBase
}

func MakeFocusLinesApp() *FocusLinesApp {
	app := &FocusLinesApp{}
	app.AppBase = apps.MakeApp(app.handleEvent, nil)
	return app
}

func (app *FocusLinesApp) Name() string {
	return "focus-lines"
}

func (app *FocusLinesApp) Start() error {
	app.SendOutput([]rofi.OutputUpdate{
		rofi.PromptUpdate{"Select an entry to focus other entry"},
		rofi.AddAllLinesUpdate{[]rofi.RofiBlocksLine{
			{Text: "focus entry 3", Data: "3"},
			{Text: "focus entry 2", Data: "2"},
			{Text: "focus entry 1000", Data: "1000"},
			{Text: "focus entry 1", Data: "1"},
			{Text: "focus entry 0", Data: "0"},
		}},
	})
	return app.AppBase.Start()
}

func (app *FocusLinesApp) handleEvent(event rofi.RofiBlocksEvent) error {
	if event.Name == rofi.SELECT_ENTRY {
		if toFocus, err := strconv.Atoi(event.Data); err != nil {
			return err
		} else {
			app.SendOutput([]rofi.OutputUpdate{
				rofi.ActiveEntryUpdate{toFocus},
			})
		}
	}
	return nil
}
