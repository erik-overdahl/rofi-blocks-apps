package rofi

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

type LineTextUpdate struct {
	Index int
	Value string
}

func (u LineTextUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= linesChanged
	line := output.Lines[u.Index]
	line.Text = u.Value
}

type LineIconUpdate struct {
	Index int
	Value string
}

func (u LineIconUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= linesChanged
	line := output.Lines[u.Index]
	line.Icon = u.Value
}

type LineDataUpdate struct {
	Index int
	Value string
}

func (u LineDataUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= linesChanged
	line := output.Lines[u.Index]
	line.Data = u.Value
}

type LineMarkupUpdate struct {
	Index int
	Value bool
}

func (u LineMarkupUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= linesChanged
	line := output.Lines[u.Index]
	line.Markup = u.Value
}

type LineHighlightUpdate struct {
	Index int
	Value bool
}

func (u LineHighlightUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= linesChanged
	line := output.Lines[u.Index]
	line.Highlight = u.Value
}

type LineUrgentUpdate struct {
	Index int
	Value bool
}

func (u LineUrgentUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= linesChanged
	line := output.Lines[u.Index]
	line.Highlight = u.Value
}

type RemoveLineUpdate struct {
	Line *RofiBlocksLine
}

func (u RemoveLineUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= linesChanged
	for i := range output.Lines {
		if output.Lines[i] == u.Line {
			output.Lines = append(output.Lines[:i], output.Lines[i+1:]...)
			break
		}
	}
}

type AddLineUpdate struct {
	Prepend bool
	Line    *RofiBlocksLine
}

func (u AddLineUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= linesChanged
	if u.Prepend {
		output.Lines = append([]*RofiBlocksLine{u.Line}, output.Lines...)
	} else {
		output.Lines = append(output.Lines, u.Line)
	}
}

type AddAllLinesUpdate struct {
	Lines []*RofiBlocksLine
}

func (u AddAllLinesUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= linesChanged
	output.Lines = append(output.Lines, u.Lines...)
}

type ReplaceLineUpdate struct {
	Existing, New *RofiBlocksLine
}

func (u ReplaceLineUpdate) Apply(output *RofiBlocksOutput) {
	output.Changes |= linesChanged
	for i, line := range output.Lines {
		if u.Existing == line {
			output.Lines[i] = u.New
		}
	}
}

type ClearLinesUpdate struct{}

func (u ClearLinesUpdate) Apply(output *RofiBlocksOutput) {
	output.Lines = []*RofiBlocksLine{}
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
	output.Lines = make([]*RofiBlocksLine, len(u.Output.Lines), len(u.Output.Lines))
	for i := range u.Output.Lines {
		output.Lines[i] = u.Output.Lines[i]
	}
	output.Changes = 0xFF
}
