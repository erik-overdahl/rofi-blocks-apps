package main

import (
	"fmt"
	"time"
)

type ShowLinesApp struct {
	updates chan []OutputUpdate
}

func (app *ShowLinesApp) Init(updates chan []OutputUpdate) {
	app.updates = updates
	currentTime := time.Now().Format("03:04:05")
	initial := []OutputUpdate{
		MessageUpdate{value: fmt.Sprintf("Updating message! \n\"Current time\" : %s", currentTime)},
		PromptUpdate{value: fmt.Sprintf("prompt %s", currentTime)},
		OverlayUpdate{value: fmt.Sprintf("Current overlay: %s", currentTime)},
		ActiveEntryUpdate{value: -1},
		SetLinesUpdate{lines: []*RofiBlocksLine{
			{Text: fmt.Sprintf("Also updates menu option text %s", currentTime)},
			{Text: "Line with urgent flag", Urgent: true},
			{Text: "Line with highlight flag", Highlight: true},
			{Text: "multi-byte unicode: •"},
			{Text: `icon unicode character: 😀`},
			{Text: "folder icon", Icon: "folder"},
			{Text: "Line <i>with</i> <b>markup</b> <b><i>flag</i></b>", Markup: true},
			{Text: "Line <i>toggling</i> <b>markup</b> flag", Markup: true},
			{Text: "Line <i>without</i> <b>markup</b> <b><i>flag</i></b>", Markup: false},
			{Text: "Line with <b><i>all</i></b> flags", Urgent: true, Highlight: true, Markup: true},
		}},
	}
	app.updates <- initial
}

func (app *ShowLinesApp) Loop() {
	toggleMarkup := true
	for {
		time.Sleep(1 * time.Second)
		currentTime := time.Now().Format("03:04:05")
		toggleMarkup = !toggleMarkup
		app.updates <- []OutputUpdate{
			MessageUpdate{value: fmt.Sprintf("Updating message! \n\"Current time\" : %s", currentTime)},
			PromptUpdate{value: fmt.Sprintf("prompt %s", currentTime)},
			OverlayUpdate{value: fmt.Sprintf("Current overlay: %s", currentTime)},
			LineTextUpdate{index: 0, value: fmt.Sprintf("Also updates menu option text %s", currentTime)},
			LineMarkupUpdate{index: 6, value: toggleMarkup},
		}
	}
}
