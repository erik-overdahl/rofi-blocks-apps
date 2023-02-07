package examples

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/erik-overdahl/rofi-blocks-apps/pkg/rofi"
	"github.com/erik-overdahl/rofi-blocks-apps/pkg/apps"
)

type VolumeApp struct {
	apps.AppBase
	currentVolume int
	muted         bool
}

func MakeVolumeApp() *VolumeApp {
	app := &VolumeApp{}
	app.AppBase = apps.MakeApp(app.handleEvent, nil)
	return app
}

func (app *VolumeApp) Start() error {
	initial := []rofi.OutputUpdate{
		rofi.PromptUpdate{"Volume"},
		rofi.AddAllLinesUpdate{Lines: []rofi.RofiBlocksLine{
			{Text: "volume up", Icon: rofi.ICONS_DIR + "/status/audio-volume-high-symbolic.svg"},
			{Text: "volume down", Icon: rofi.ICONS_DIR + "/status/audio-volume-low-symbolic.svg"},
			{Text: "toggle mute", Icon: rofi.ICONS_DIR + "/status/audio-volume-muted-symbolic.svg"},
		}},
	}
	app.SendOutput(initial)
	go app.listenVolumeChange()
	return app.AppBase.Start()
}

func (app *VolumeApp) handleEvent(event rofi.RofiBlocksEvent) error {
	if event.Name == rofi.SELECT_ENTRY {
		switch event.Value {
		case "volume up":
			exec.Command("pactl", "set-sink-volume", "@DEFAULT_SINK@", "+5%").Run()
		case "volume down":
			exec.Command("pactl", "set-sink-volume", "@DEFAULT_SINK@", "-5%").Run()
		case "toggle mute":
			exec.Command("pactl", "set-sink-mute", "@DEFAULT_SINK@", "toggle").Run()
		}
	}
	return nil
}

func (app *VolumeApp) listenVolumeChange() {
	var message string
	for {
		select {
		case <-app.Done():
			return
		default:
			muteOutput, err := exec.Command("pactl", "get-sink-mute", "@DEFAULT_SINK@").Output()
			if err != nil {
				log.Printf("Failed to check mute: %v\n", err)
			}
			if string(muteOutput) == "Mute: yes" {
				app.muted = true
			} else {
				volumeOutput, err := exec.Command("pactl", "get-sink-volume", "@DEFAULT_SINK@").Output()
				if err != nil {
					log.Printf("Failed to check volume: %v\n", err)
				}
				pieces := strings.Split(string(volumeOutput), " ")
				volumeStr := strings.TrimRight(pieces[5], "%")
				app.currentVolume, _ = strconv.Atoi(volumeStr)
			}
			if app.muted {
				message = "Current volume: muted"
			} else {
				message = fmt.Sprintf("Current volume: %d%%", app.currentVolume)
			}
			app.SendOutput([]rofi.OutputUpdate{
				rofi.MessageUpdate{Value: message},
			})
			time.Sleep(500 * time.Millisecond)
		}
	}
}
