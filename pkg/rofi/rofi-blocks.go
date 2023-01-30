package rofi

import (
	"encoding/json"
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

var ICONS_DIR = "/usr/share/icons/Adwaita/scalable"

type RofiBlocksEvent struct {
	Name  string           `json:"name"`
	Value string           `json:"value"`
	Data  string           `json:"data"`
	Prev  *RofiBlocksEvent `json:"-"`
}

type RofiBlocksLine struct {
	Text      string `json:"text,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Data      string `json:"data,omitempty"`
	Urgent    bool   `json:"urgent,omitempty"`
	Highlight bool   `json:"highlight,omitempty"`
	Markup    bool   `json:"markup,omitempty"`
}

// for the output struct, we send just the delta (with the exception that,
// if any of the lines changed, we have to send all of them again)
type RofiBlocksOutput struct {
	Changes     uint8
	Prompt      string
	Message     string
	Overlay     string
	Input       string
	ActiveEntry int
	InputAction string
	Lines       []*RofiBlocksLine
}

func MakeRofiBlocksOutput() *RofiBlocksOutput {
	return &RofiBlocksOutput{
		InputAction: INPUT_ACTION_FILTER, // default
		Lines: []*RofiBlocksLine{},
	}
}

func (o *RofiBlocksOutput) ChangeAll() {
	o.Changes = 0xFF
}

func (this *RofiBlocksOutput) ApplyAll(updates []OutputUpdate) {
	for _, change := range updates {
		change.Apply(this)
	}
}

func (o RofiBlocksOutput) MarshalJson() ([]byte, error) {
	pieces := map[string]any{}
	if o.Changes&promptChanged > 0 {
		pieces["prompt"] = o.Prompt
	}
	if o.Changes&messageChanged > 0 {
		pieces["message"] = o.Message
	}
	if o.Changes&overlayChanged > 0 {
		pieces["overlay"] = o.Overlay
	}
	if o.Changes&inputChanged > 0 {
		pieces["input"] = o.Input
	}
	if o.Changes&activeEntryChanged > 0 {
		pieces["active entry"] = o.ActiveEntry
	}
	if o.Changes&inputActionChanged > 0 {
		pieces["input action"] = o.InputAction
	}
	pieces["lines"] = o.Lines
	return json.Marshal(pieces)
}

