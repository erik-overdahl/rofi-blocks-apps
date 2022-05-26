package main

type OutputUpdate interface {
	Apply(*RofiBlocksOutput)
}

type PromptUpdate struct {
	value string
}

func (u PromptUpdate) Apply(output *RofiBlocksOutput) {
	output.ChangePrompt = true
	output.Prompt = u.value
}

type MessageUpdate struct {
	value string
}

func (u MessageUpdate) Apply(output *RofiBlocksOutput) {
	output.ChangeMessage = true
	output.Message = u.value
}

type OverlayUpdate struct {
	value string
}

func (u OverlayUpdate) Apply(output *RofiBlocksOutput) {
	output.ChangeOverlay = true
	output.Overlay = u.value
}

type InputUpdate struct {
	value string
}

func (u InputUpdate) Apply(output *RofiBlocksOutput) {
	output.ChangeInput = true
	output.Input = u.value
}

type InputActionUpdate struct {
	value string
}

func (u InputActionUpdate) Apply(output *RofiBlocksOutput) {
	output.ChangeInputAction = true
	output.InputAction = u.value
}

type ActiveEntryUpdate struct {
	value int
}

func (u ActiveEntryUpdate) Apply(output *RofiBlocksOutput) {
	output.ChangeActiveEntry = true
	output.ActiveEntry = u.value
}

type LineTextUpdate struct {
	index int
	value string
}

func (u LineTextUpdate) Apply(output *RofiBlocksOutput) {
	output.ChangeLines = true
	line := output.Lines[u.index]
	line.Text = u.value
}

type LineIconUpdate struct {
	index int
	value string
}

func (u LineIconUpdate) Apply(output *RofiBlocksOutput) {
	output.ChangeLines = true
	line := output.Lines[u.index]
	line.Icon = u.value
}

type LineDataUpdate struct {
	index int
	value string
}

func (u LineDataUpdate) Apply(output *RofiBlocksOutput) {
	output.ChangeLines = true
	line := output.Lines[u.index]
	line.Data = u.value
}

type LineMarkupUpdate struct {
	index int
	value bool
}

func (u LineMarkupUpdate) Apply(output *RofiBlocksOutput) {
	output.ChangeLines = true
	line := output.Lines[u.index]
	line.Markup = u.value
}

type LineHighlightUpdate struct {
	index int
	value bool
}

func (u LineHighlightUpdate) Apply(output *RofiBlocksOutput) {
	output.ChangeLines = true
	line := output.Lines[u.index]
	line.Highlight = u.value
}

type LineUrgentUpdate struct {
	index int
	value bool
}

func (u LineUrgentUpdate) Apply(output *RofiBlocksOutput) {
	output.ChangeLines = true
	line := output.Lines[u.index]
	line.Highlight = u.value
}

type RemoveLineUpdate struct {
	line *RofiBlocksLine
}

func (u RemoveLineUpdate) Apply(output *RofiBlocksOutput) {
	output.ChangeLines = true
	for i, l := range output.Lines {
		if l == u.line {
			output.Lines = append(output.Lines[:i], output.Lines[i+1:]...)
			break
		}
	}
}

type AddLineUpdate struct {
	prepend bool
	line    *RofiBlocksLine
}

func (u AddLineUpdate) Apply(output *RofiBlocksOutput) {
	output.ChangeLines = true
	if u.prepend {
		output.Lines = append([]*RofiBlocksLine{u.line}, output.Lines...)
	} else {
		output.Lines = append(output.Lines, u.line)
	}
}

type ReplaceLineUpdate struct {
	existing, new *RofiBlocksLine
}

func (u ReplaceLineUpdate) Apply(output *RofiBlocksOutput) {
	output.ChangeLines = true
	for i, line := range output.Lines {
		if u.existing == line {
			output.Lines[i] = u.new
		}
	}
}

type ClearLinesUpdate struct{}

func (u ClearLinesUpdate) Apply(output *RofiBlocksOutput) {
	output.ChangeLines = true
	output.Lines = []*RofiBlocksLine{}
}

type SetLinesUpdate struct {
	lines []*RofiBlocksLine
}

func (u SetLinesUpdate) Apply(output *RofiBlocksOutput) {
	output.ChangeLines = true
	output.Lines = u.lines
}

type SnapshotState struct {
	Output *RofiBlocksOutput
}

func (u SnapshotState) Apply(output *RofiBlocksOutput) {
	u.Output.Prompt = output.Prompt
	u.Output.Message = output.Message
	u.Output.Overlay = output.Overlay
	u.Output.Input = output.Input
	u.Output.InputAction = output.InputAction
	u.Output.ActiveEntry = output.ActiveEntry
	u.Output.Lines = make([]*RofiBlocksLine, len(output.Lines), len(output.Lines))

	for i, line := range output.Lines {
		u.Output.Lines[i] = line
	}
	logger.Println("Snapshotting state")
}

type RestoreState struct {
	Output *RofiBlocksOutput
}

func (u RestoreState) Apply(output *RofiBlocksOutput) {
	output.Prompt = u.Output.Prompt
	output.ChangePrompt = true
	output.Message = u.Output.Message
	output.ChangeMessage = true
	output.Overlay = u.Output.Overlay
	output.ChangeOverlay = true
	output.Input = u.Output.Input
	output.ChangeInput = true
	output.InputAction = u.Output.InputAction
	output.ChangeInputAction = true
	output.ActiveEntry = u.Output.ActiveEntry
	output.ChangeActiveEntry = true
	output.Lines = make([]*RofiBlocksLine, len(u.Output.Lines), len(u.Output.Lines))
	for i, line := range u.Output.Lines {
		output.Lines[i] = line
	}
	output.ChangeLines = true
	logger.Println("Restoring state")
	// logger.Printf("Restoring to:\n%+v\n", output.ToString())
}
