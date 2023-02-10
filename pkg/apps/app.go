package apps

import (
	"context"
	"log"

	"github.com/erik-overdahl/rofi-blocks-apps/pkg/rofi"
)

type App interface {
	Start() error
	Stop() error
	Done() <-chan struct{}
	Foreground() *rofi.RofiBlocksOutput
	Background()
	Input() chan<- rofi.Event
	Output() <-chan []rofi.OutputUpdate
	ShouldReceiveInBackground() bool
}

type AppBase struct {
	inForeground bool
	running 	 bool
	state        *rofi.RofiBlocksOutput
	input        chan rofi.Event
	output       chan []rofi.OutputUpdate
	ctx          context.Context
	ctxCancel    context.CancelFunc
	handleEvent  func(rofi.Event) error
	loop         func()
}

func MakeApp(handleEvent func(rofi.Event) error, loop func()) AppBase {
	return AppBase{
		state:       rofi.MakeRofiBlocksOutput(),
		input:       make(chan rofi.Event),
		output:      make(chan []rofi.OutputUpdate),
		handleEvent: handleEvent,
		loop:        loop,
	}
}

func (app *AppBase) Start() error {
	app.ctx, app.ctxCancel = context.WithCancel(context.Background())
	if app.loop != nil {
		go app.loop()
	}
	if app.handleEvent == nil {
		app.handleEvent = func(event rofi.Event) error { return nil; }
	}
	go func() {
		for {
			select {
			case <-app.ctx.Done():
				return
			case event, ok := <-app.input:
				if !ok {
					return
				}
				if err := app.handleEvent(event); err != nil {
					log.Printf("Error handling event %v: %v", event, err)
				}
			}
		}
	}()
	app.running = true
	return nil
}

func (app *AppBase) Stop() error {
	app.ctxCancel()
	app.running = false
	return nil
}

func (app *AppBase) Done() <-chan struct{} {
	return app.ctx.Done()
}

func (app *AppBase) Background() {
	app.inForeground = false
}

func (app *AppBase) Foreground() *rofi.RofiBlocksOutput {
	app.inForeground = true
	return app.state
}

func (app *AppBase) Input() chan<- rofi.Event {
	return app.input
}

func (app *AppBase) Output() <-chan []rofi.OutputUpdate {
	return app.output
}

func (app *AppBase) ShouldReceiveInBackground() bool {
	return false
}

func (app *AppBase) SendOutput(updates []rofi.OutputUpdate) {
	app.state.ApplyAll(updates)
	if app.inForeground {
		app.output <- updates
	}
}
