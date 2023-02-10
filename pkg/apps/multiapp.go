package apps

import (
	"fmt"

	"github.com/erik-overdahl/rofi-blocks-apps/pkg/rofi"
)

type MultiApp struct {
	AppBase
	swap    chan App
	apps    []App
	current App
}

func MakeMultiApp() *MultiApp {
	a := &MultiApp{
		swap: make(chan App),
		apps: []App{},
	}
	a.AppBase = MakeApp(a.handleEvent, a.loop)
	a.inForeground = true
	return a
}

func (this *MultiApp) Start() error {
	for i, app := range this.apps {
		if err := app.Start(); err != nil {
			for _, a := range this.apps[:i] {
				a.Stop()
			}
			return fmt.Errorf("Multiapp failed to start: error in app %d: %v", i, err)
		}
	}
	if err := this.AppBase.Start(); err != nil {
		return fmt.Errorf("Multiapp failed to start: %v", err)
	}
	if len(this.apps) > 0 {
		this.current = this.apps[0]
		state := this.current.Foreground()
		this.SendOutput([]rofi.OutputUpdate{
			rofi.RestoreState{Output: state},
		})
	}
	return nil
}

func (this *MultiApp) Add(apps ...App) {
	this.apps = append(this.apps, apps...)
}

func (this *MultiApp) SwitchTo(app int) error {
	if app < 0 || len(this.apps) <= app || len(this.apps) == 0 {
		return fmt.Errorf("App %d does not exist", app)
	}
	if this.running {
		this.swap <- this.apps[app]
	} else {
		this.current.Background()
		this.current = this.apps[app]
	}
	return nil
}

func (this *MultiApp) handleEvent(event rofi.Event) error {
	// TODO: currently no CUSTOM_KEY events are handed down to child apps
	switch event := event.(type) {
	case *rofi.CustomKeyEvent:
		return this.SwitchTo(event.KeyId - 1)
	default:
		for _, app := range this.apps {
			if app == this.current || app.ShouldReceiveInBackground() {
				// it would be nice if this were async, but hopefully
				// won't be an issue
				//
				// just starting a goroutine for each app would
				// potentially send events out of order, so we would
				// need to queue them instead
				app.Input() <- event
			}
		}
	}
	return nil
}

func (this *MultiApp) loop() {
	for {
		select {
		case new := <-this.swap:
			this.current.Background()
			this.current = new
			state := this.current.Foreground()
			this.SendOutput([]rofi.OutputUpdate{
				rofi.RestoreState{Output: state},
			})
		case update := <-this.current.Output():
			this.SendOutput(update)
		}
	}
}
