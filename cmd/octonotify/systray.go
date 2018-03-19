package main

import (
	"os"

	"github.com/getlantern/systray"
	"github.com/sirupsen/logrus"
	"github.com/twexler/octonotify/icons"
)

func initSystray() {
	systray.Run(onReady, onExit)
}

func onReady() {
	data, err := icons.Asset(icons.IconPath)
	if err != nil {
		logrus.WithError(err).Error("unable to load icons")
		return
	}
	systray.SetIcon(data)
	systray.SetTooltip("octonotify")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuit.ClickedCh
		os.Exit(0)
	}()
}

func onExit() {

}
