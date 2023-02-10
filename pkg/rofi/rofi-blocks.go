package rofi

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"
)

const (
	INPUT_ACTION_FILTER = "filter"
	INPUT_ACTION_SEND   = "send"
)

var ICONS_DIR = "/usr/share/icons/Adwaita/scalable"

var (
	lineId      LineId = 0
	lineIdMutex sync.Mutex
)

type LineId int

func NewLineId() LineId {
	lineIdMutex.Lock()
	defer lineIdMutex.Unlock()
	lineId++
	return lineId
}

type RofiBlocksLine struct {
	Id        LineId
	Text      string
	Icon      string
	Urgent    bool
	Highlight bool
	Markup    bool
}

type rawLine struct {
	Data      string `json:"data,omitempty"`
	Text      string `json:"text,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Urgent    bool   `json:"urgent,omitempty"`
	Highlight bool   `json:"highlight,omitempty"`
	Markup    bool   `json:"markup,omitempty"`
}

var xmlEscapes = []struct{ old, new string }{
	{"&", "&amp;"},
	{`"`, "&quot;"},
	{"'", "&apos;"},
	{"<", "&lt;"},
	{">", "&gt;"},
}

func (line RofiBlocksLine) MarshalJSON() ([]byte, error) {
	text := line.Text
	for _, change := range xmlEscapes {
		text = strings.ReplaceAll(text, change.old, change.new)
	}
	return json.Marshal(rawLine{
		Data:      strconv.Itoa(int(line.Id)),
		Text:      text,
		Icon:      line.Icon,
		Urgent:    line.Urgent,
		Highlight: line.Highlight,
		Markup:    line.Markup,
	})
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
	Lines       []RofiBlocksLine
}

func MakeRofiBlocksOutput() *RofiBlocksOutput {
	return &RofiBlocksOutput{
		InputAction: INPUT_ACTION_FILTER, // default
		Lines:       []RofiBlocksLine{},
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
