package main

type App interface {
	Init(chan []OutputUpdate)
}

type EventHandlingApp interface {
	App
	// if an app wants to change the output to the existing rofi
	// process, it should send on the outputUpdates channel.
	// if it wants want change the existing process, say to change
	// the theme, it should return the args it wants passed in a
	// _processChange struct. otherwise this should return nil
	HandleEvent(RofiBlocksEvent) *ProcessChange
}

type LoopingApp interface {
	App
	Loop()
}
