package main

import (
	"fmt"
)

type AppNotFound struct {
	command string
}

func (app *AppNotFound) Init(updates chan []OutputUpdate) {
	updates <- []OutputUpdate{
		MessageUpdate{value: fmt.Sprintf("Unknown command '%s'", app.command)},
		AddLineUpdate{line: &RofiBlocksLine{Text: "ok"}},
	}
}

func (app *AppNotFound) RofiArgs() []string {
	return nil
}
