package rofi

type Event interface {
	Prev() Event
}

type rawEvent struct {
	Name  string `json:"name"`
	Value string `json:"Value"`
	Data  string `json:"Data"`
}

type InputChangeEvent struct {
	Value string
	prev  Event
}

func (event *InputChangeEvent) Prev() Event {
	return event.prev
}

type CustomKeyEvent struct {
	KeyId int // 1 - 19
	prev  Event
}

func (event *CustomKeyEvent) Prev() Event {
	return event.prev
}

type ActiveEntryEvent struct {
	LineId int
	prev  Event
}

func (event *ActiveEntryEvent) Prev() Event {
	return event.prev
}

type SelectEntryEvent struct {
	LineId int
	prev   Event
}

func (event *SelectEntryEvent) Prev() Event {
	return event.prev
}

type DeleteEntryEvent struct {
	LineId int
	prev   Event
}

func (event *DeleteEntryEvent) Prev() Event {
	return event.prev
}

type ExecCustomEntryEvent struct {
	Value string
	prev  Event
}

func (event *ExecCustomEntryEvent) Prev() Event {
	return event.prev
}
