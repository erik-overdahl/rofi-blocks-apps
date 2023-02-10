package examples

import (
	"github.com/erik-overdahl/rofi-blocks-apps/pkg/rofi"
	"github.com/erik-overdahl/rofi-blocks-apps/pkg/apps"
)

type FocusLinesApp struct {
	apps.AppBase
	lines []*rofi.RofiBlocksLine
	focus []int
}

func MakeFocusLinesApp() *FocusLinesApp {
	app := &FocusLinesApp{
		lines: []*rofi.RofiBlocksLine{
			{Text: "focus entry 3", Id: rofi.NewLineId()},
			{Text: "focus entry 2", Id: rofi.NewLineId()},
			{Text: "focus entry 1000", Id: rofi.NewLineId()},
			{Text: "focus entry 1", Id: rofi.NewLineId()},
			{Text: "focus entry 0", Id: rofi.NewLineId()},
		},
		focus: []int{3, 2, 1000, 1, 0},
	}
	app.AppBase = apps.MakeApp(app.handleEvent, nil)
	return app
}

func (app *FocusLinesApp) Name() string {
	return "focus-lines"
}

func (app *FocusLinesApp) Start() error {
	lines := make([]rofi.RofiBlocksLine, len(app.lines), len(app.lines))
	for i, line := range app.lines {
		lines[i] = *line
	}
	app.SendOutput([]rofi.OutputUpdate{
		rofi.PromptUpdate{"Select an entry to focus other entry"},
		rofi.AddAllLinesUpdate{Lines: lines},
	})
	return app.AppBase.Start()
}

func (app *FocusLinesApp) handleEvent(event rofi.Event) error {
	switch event := event.(type) {
	case *rofi.SelectEntryEvent:
		for i, l := range app.lines {
			if l.Id == event.LineId {
				app.SendOutput([]rofi.OutputUpdate{
					rofi.ActiveEntryUpdate{app.focus[i]},
				})
			}
		}
	}
	return nil
}
