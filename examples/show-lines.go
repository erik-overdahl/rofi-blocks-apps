package examples

import (
	"fmt"
	"time"

	"github.com/erik-overdahl/rofi-blocks-apps/pkg/rofi"
	"github.com/erik-overdahl/rofi-blocks-apps/pkg/apps"
)

type ShowLinesApp struct {
	apps.AppBase
	lines []*rofi.RofiBlocksLine
}

func MakeShowLinesApp() *ShowLinesApp {
	app := &ShowLinesApp{}
	app.AppBase = apps.MakeApp(nil, app.loop)
	return app
}

func (app *ShowLinesApp) Start() error {
	currentTime := time.Now().Format("03:04:05")
	initial := []rofi.OutputUpdate{
		rofi.MessageUpdate{fmt.Sprintf("Updating message! \n\"Current time\" : %s", currentTime)},
		rofi.PromptUpdate{fmt.Sprintf("prompt %s", currentTime)},
		rofi.OverlayUpdate{fmt.Sprintf("Current overlay: %s", currentTime)},
		rofi.AddAllLinesUpdate{[]*rofi.RofiBlocksLine{
			{Text: fmt.Sprintf("Also updates menu option text %s", currentTime)},
			{Text: "Line with urgent flag", Urgent: true},
			{Text: "Line with highlight flag", Highlight: true},
			{Text: "multi-byte unicode: •"},
			{Text: `icon unicode character: 😀`},
			{Text: "folder icon", Icon: "folder"},
			{Text: "Line <i>with</i> <b>markup</b> <b><i>flag</i></b>", Markup: true},
			{Text: "Line <i>toggling</i> <b>markup</b> flag", Markup: true},
			{Text: "Line <i>without</i> <b>markup</b> <b><i>flag</i></b>", Markup: false},
			{Text: "Line with <b><i>all</i></b> flags", Urgent: true, Highlight: true, Markup: true}}},
	}
	app.SendOutput(initial)
	return app.AppBase.Start()
}

func (app *ShowLinesApp) loop() {
	toggleMarkup := true
	for {
		time.Sleep(1 * time.Second)
		currentTime := time.Now().Format("03:04:05")
		toggleMarkup = !toggleMarkup
		app.SendOutput([]rofi.OutputUpdate{
			rofi.MessageUpdate{fmt.Sprintf("Updating message! \n\"Current time\" : %s", currentTime)},
			rofi.PromptUpdate{fmt.Sprintf("prompt %s", currentTime)},
			rofi.OverlayUpdate{fmt.Sprintf("Current overlay: %s", currentTime)},
			rofi.LineTextUpdate{Index: 0, Value: fmt.Sprintf("Also updates menu option text %s", currentTime)},
			rofi.LineMarkupUpdate{Index: 7, Value: toggleMarkup},
		})
	}
}
