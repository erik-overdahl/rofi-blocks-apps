package main

import "strconv"

type FocusLinesApp struct {
	updates chan []OutputUpdate
}

func (app *FocusLinesApp) Name() string {
	return "focus lines"
}

func (app *FocusLinesApp) RofiArgs() []string {
	return nil
}

func (app *FocusLinesApp) Init(updates chan []OutputUpdate) {
	logger.Println("Initializing Focus Lines app")
	logger.Println("Sending initial ouput")
	updates <- []OutputUpdate{
		PromptUpdate{value: "Select an entry to focus other entry"},
		SetLinesUpdate{lines: []*RofiBlocksLine{
			{Text: "focus entry 3", Data: "3"},
			{Text: "focus entry 2", Data: "2"},
			{Text: "focus entry 1000", Data: "1000"},
			{Text: "focus entry 1", Data: "1"},
			{Text: "focus entry 0", Data: "0"},
		}},
	}
	app.updates = updates
}

func (app *FocusLinesApp) HandleEvent(event RofiBlocksEvent) *ProcessChange {
	if event.Name == SELECT_ENTRY {
		toFocus, _ := strconv.Atoi(event.Data)
		app.updates <- []OutputUpdate{
			ActiveEntryUpdate{value: toFocus},
			InputUpdate{value: ""},
		}
	}
	return nil
}
