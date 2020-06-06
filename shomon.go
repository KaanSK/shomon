package main

import (
	"os"

	"github.com/KaanSK/shomon/pkg/conf"
	lw "github.com/KaanSK/shomon/pkg/logwrapper"
	"github.com/KaanSK/shomon/pkg/shodan"
	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

func init() {
	parser := flags.NewParser(&conf.Config, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}
	if conf.Config.Verbose {
		lw.Logger.Formatter = &logrus.JSONFormatter{}
		lw.Logger.SetReportCaller(true)
		lw.Logger.SetLevel(logrus.DebugLevel)
	}
}

func neverExit() {
	defer func() {
		if err := recover(); err != nil {
			lw.Logger.Error(err)
			go neverExit()
		}
	}()
	err := shodan.ListenAlerts()
	if err != nil {
		lw.Logger.Error(err)
		go neverExit()
	}
}

func main() {
	lw.Logger.Info("main process started")
	go neverExit()
	select {}
}
