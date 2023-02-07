package rofi

import "github.com/google/uuid"

const (
	promptChanged      = 0x01
	messageChanged     = 0x02
	overlayChanged     = 0x04
	inputChanged       = 0x08
	activeEntryChanged = 0x10
	inputActionChanged = 0x20
	linesChanged 	   = 0x40
)

type OutputUpdate interface {
	Apply(*RofiBlocksOutput)
}

type PromptUpdate struct {
	Value string
}

func (u PromptUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= promptChanged
	output.Prompt = u.Value
}

type MessageUpdate struct {
	Value string
}

func (u MessageUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= messageChanged
	output.Message = u.Value
}

type OverlayUpdate struct {
	Value string
}

func (u OverlayUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= overlayChanged
	output.Overlay = u.Value
}

type InputUpdate struct {
	Value string
}

func (u InputUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= inputChanged
	output.Input = u.Value
}

type InputActionUpdate struct {
	Value string
}

func (u InputActionUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= inputActionChanged
	output.InputAction = u.Value
}

type ActiveEntryUpdate struct {
	Value int
}

func (u ActiveEntryUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= activeEntryChanged
	output.ActiveEntry = u.Value
}

type LineUpdate struct {
	Line RofiBlocksLine
}

func (u LineUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= linesChanged
	for i, line := range output.Lines {
		if line.Id == u.Line.Id {
			output.Lines[i] = u.Line
			break
		}
	}
}

type RemoveLineUpdate struct {
	LineId uuid.UUID
}

func (u RemoveLineUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= linesChanged
	for i, line := range output.Lines {
		if line.Id == u.LineId {
			output.Lines = append(output.Lines[:i], output.Lines[i+1:]...)
			break
		}
	}
}

type AddLineUpdate struct {
	Prepend bool
	Line    RofiBlocksLine
}

func (u AddLineUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= linesChanged
	if u.Prepend {
		output.Lines = append([]RofiBlocksLine{u.Line}, output.Lines...)
	} else {
		output.Lines = append(output.Lines, u.Line)
	}
}

type AddAllLinesUpdate struct {
	Lines []RofiBlocksLine
}

func (u AddAllLinesUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= linesChanged
	output.Lines = append(output.Lines, u.Lines...)
}

type ClearLinesUpdate struct{}

func (u ClearLinesUpdate) Apply(output *RofiBlocksOutput) {
	output.Lines = []RofiBlocksLine{}
	output.Changes |= linesChanged
}

type RestoreState struct {
	Output *RofiBlocksOutput
}

func (u RestoreState) Apply(output *RofiBlocksOutput) {
	if u.Output == output {
		return
	}
	output.Prompt = u.Output.Prompt
	output.Message = u.Output.Message
	output.Overlay = u.Output.Overlay
	output.Input = u.Output.Input
	output.InputAction = u.Output.InputAction
	output.ActiveEntry = u.Output.ActiveEntry
	output.Lines = make([]RofiBlocksLine, len(u.Output.Lines), len(u.Output.Lines))
	for i := range u.Output.Lines {
		output.Lines[i] = u.Output.Lines[i]
	}
	output.Changes = 0xFF
}
