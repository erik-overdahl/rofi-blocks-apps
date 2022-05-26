package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	INPUT_CHANGE      = "input change"
	CUSTOM_KEY        = "custom key"
	ACTIVE_ENTRY      = "active entry"
	SELECT_ENTRY      = "select entry"
	DELETE_ENTRY      = "delete entry"
	EXEC_CUSTOM_INPUT = "execute custom input"

	INPUT_ACTION_FILTER = "filter"
	INPUT_ACTION_SEND   = "send"
)

type RofiBlocksEvent struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Data  string `json:"data"`
	Prev  *RofiBlocksEvent
}

type RofiBlocksLine struct {
	Text      string
	Icon      string
	Data      string
	Urgent    bool
	Highlight bool
	Markup    bool
}

func (l RofiBlocksLine) ToString() string {
	pieces := []string{fmt.Sprintf(`"urgent":%t, "highlight":%t, "markup":%t`, l.Urgent, l.Highlight, l.Markup)}
	if l.Text != "" {
		pieces = append(pieces, fmt.Sprintf(`"text":"%s"`, jsonEscape(l.Text)))
	}
	if l.Icon != "" {
		pieces = append(pieces, fmt.Sprintf(`"icon":"%s"`, l.Icon))
	}
	if l.Data != "" {
		pieces = append(pieces, fmt.Sprintf(`"data":"%s"`, l.Data))
	}
	return fmt.Sprintf("{%s}", strings.Join(pieces, ","))
}

type RofiBlocksOutput struct {
	Prompt            string
	ChangePrompt      bool
	Message           string
	ChangeMessage     bool
	Overlay           string
	ChangeOverlay     bool
	Input             string
	ChangeInput       bool
	ActiveEntry       int
	ChangeActiveEntry bool
	InputAction       string
	ChangeInputAction bool
	Lines             []*RofiBlocksLine
	ChangeLines       bool
}

func (o *RofiBlocksOutput) ResetFlags() {
	o.ChangePrompt = false
	o.ChangeMessage = false
	o.ChangeOverlay = false
	o.ChangeInput = false
	o.ChangeActiveEntry = false
	o.ChangeInputAction = false
	o.ChangeLines = false
}

func (o RofiBlocksOutput) ToString() string {
	pieces := []string{}
	if o.ChangePrompt {
		pieces = append(pieces, fmt.Sprintf(`"prompt":"%s"`, jsonEscape(o.Prompt)))
	}
	if o.ChangeMessage {
		pieces = append(pieces, fmt.Sprintf(`"message":"%s"`, jsonEscape(o.Message)))
	}
	if o.ChangeOverlay {
		pieces = append(pieces, fmt.Sprintf(`"overlay":"%s"`, jsonEscape(o.Overlay)))
	}
	if o.ChangeInput {
		pieces = append(pieces, fmt.Sprintf(`"input":"%s"`, jsonEscape(o.Input)))
	}
	if o.ChangeActiveEntry {
		pieces = append(pieces, fmt.Sprintf(`"active entry": %d`, o.ActiveEntry))
	}
	if o.ChangeInputAction {
		pieces = append(pieces, fmt.Sprintf(`"input action":"%s"`, o.InputAction))
	}
	if o.ChangeLines {
		lines := []string{}
		for _, l := range o.Lines {
			lines = append(lines, l.ToString())
		}
		pieces = append(pieces, fmt.Sprintf(`"lines":[%s]`, strings.Join(lines, ",")))
	}
	return fmt.Sprintf("{%s}\n", strings.Join(pieces, ","))
}

func jsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	// Trim the beginning and trailing " character
	return string(b[1 : len(b)-1])
}
