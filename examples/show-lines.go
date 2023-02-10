package examples

import (
	"fmt"
	"time"

	"github.com/erik-overdahl/rofi-blocks-apps/pkg/apps"
	"github.com/erik-overdahl/rofi-blocks-apps/pkg/rofi"
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

func (app *ShowLinesApp) Name() string {
	return "show-lines"
}

func (app *ShowLinesApp) Start() error {
	currentTime := time.Now().Format("03:04:05")
	app.lines = []*rofi.RofiBlocksLine{
		{
			Id: rofi.NewLineId(),
			Text: fmt.Sprintf("Also updates menu option text %s", currentTime),
		},
		{
			Id: rofi.NewLineId(),
			Text: "Line with urgent flag",
			Urgent: true,
		},
		{
			Id: rofi.NewLineId(),
			Text: "Line with highlight flag",
			Highlight: true,
		},
		{
			Id: rofi.NewLineId(),
			Text: "multi-byte unicode: •",
		},
		{
			Id: rofi.NewLineId(),
			Text: `icon unicode character: 😀`,
		},
		{
			Id: rofi.NewLineId(),
			Text: "folder icon",
			Icon: "folder",
		},
		{
			Id: rofi.NewLineId(),
			Text: "Line <i>with</i> <b>markup</b> <b><i>flag</i></b>",
			Markup: true,
		},
		{
			Id: rofi.NewLineId(),
			Text: "Line <i>toggling</i> <b>markup</b> flag",
			Markup: true,
		},
		{
			Id: rofi.NewLineId(),
			Text: "Line <i>without</i> <b>markup</b> <b><i>flag</i></b>",
			Markup: false,
		},
		{
			Id: rofi.NewLineId(),
			Text: "Line with <b><i>all</i></b> flags",
			Urgent: true,
			Highlight: true,
			Markup: true,
		},
	}
	addLines := rofi.AddAllLinesUpdate{Lines: make([]rofi.RofiBlocksLine, len(app.lines), len(app.lines))}
	for i, line := range app.lines {
		addLines.Lines[i] = *line
	}
	initial := []rofi.OutputUpdate{
		rofi.MessageUpdate{fmt.Sprintf("Updating message! \n\"Current time\" : %s", currentTime)},
		rofi.PromptUpdate{fmt.Sprintf("prompt %s", currentTime)},
		rofi.OverlayUpdate{fmt.Sprintf("Current overlay: %s", currentTime)},
		addLines,
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
		app.lines[0].Text = fmt.Sprintf("Also updates menu option text %s", currentTime)
		app.lines[7].Markup = toggleMarkup
		app.SendOutput([]rofi.OutputUpdate{
			rofi.MessageUpdate{fmt.Sprintf("Updating message! \n\"Current time\" : %s", currentTime)},
			rofi.PromptUpdate{fmt.Sprintf("prompt %s", currentTime)},
			rofi.OverlayUpdate{fmt.Sprintf("Current overlay: %s", currentTime)},
			rofi.LineUpdate{Line: *(app.lines[0])},
			rofi.LineUpdate{Line: *(app.lines[7])},
		})
	}
}
