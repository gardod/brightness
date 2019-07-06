package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gardod/brightness/sysmonitor"
	"github.com/gardod/brightness/systray"
)

var (
	set100 *systray.MenuItem
	set90  *systray.MenuItem
	set80  *systray.MenuItem
	set70  *systray.MenuItem
	set60  *systray.MenuItem
	set50  *systray.MenuItem
	set40  *systray.MenuItem
	set30  *systray.MenuItem
	set20  *systray.MenuItem
	set10  *systray.MenuItem
	set0   *systray.MenuItem
	exit   *systray.MenuItem

	monitors []sysmonitor.Monitor
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(getIcon("brightness.ico"))
	systray.SetTitle("Brightness")
	systray.SetTooltip("Brightness")

	setupMonitors()
	setupControls()

	listen()
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		log.Print(err)
	}
	return b
}

func setupControls() {
	set100 = systray.AddMenuItem("100%", "Set brightness")
	set90 = systray.AddMenuItem("90%", "Set brightness")
	set80 = systray.AddMenuItem("80%", "Set brightness")
	set70 = systray.AddMenuItem("70%", "Set brightness")
	set60 = systray.AddMenuItem("60%", "Set brightness")
	set50 = systray.AddMenuItem("50%", "Set brightness")
	set40 = systray.AddMenuItem("40%", "Set brightness")
	set30 = systray.AddMenuItem("30%", "Set brightness")
	set20 = systray.AddMenuItem("20%", "Set brightness")
	set10 = systray.AddMenuItem("10%", "Set brightness")
	set0 = systray.AddMenuItem("0%", "Set brightness")
	systray.AddSeparator()
	exit = systray.AddMenuItem("Exit", "Close application")
}

func setupMonitors() {
	var err error
	monitors, err = sysmonitor.GetMonitors()
	if err != nil {
		log.Print(err)
	}

	for _, monitor := range monitors {
		systray.AddMenuItem(monitor.Name, "").Disable()
	}
	systray.AddSeparator()
}

func listen() {
	for {
		select {
		case <-set100.ClickedCh:
			setBrightness(100)
		case <-set90.ClickedCh:
			setBrightness(90)
		case <-set80.ClickedCh:
			setBrightness(80)
		case <-set70.ClickedCh:
			setBrightness(70)
		case <-set60.ClickedCh:
			setBrightness(60)
		case <-set50.ClickedCh:
			setBrightness(50)
		case <-set40.ClickedCh:
			setBrightness(40)
		case <-set30.ClickedCh:
			setBrightness(30)
		case <-set20.ClickedCh:
			setBrightness(20)
		case <-set10.ClickedCh:
			setBrightness(10)
		case <-set0.ClickedCh:
			setBrightness(0)
		case <-exit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func setBrightness(perc int64) {
	for _, monitor := range monitors {
		monitor.SetBrightness(perc)
	}
	systray.SetTooltip(fmt.Sprintf("Brightness: %d%%", perc))
}

func onExit() {
	for _, monitor := range monitors {
		monitor.Destroy()
	}
}
