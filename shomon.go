package main

import (
	"os"

	"github.com/KaanSK/shomon/pkg/conf"
	lw "github.com/KaanSK/shomon/pkg/logwrapper"
	"github.com/KaanSK/shomon/pkg/shodan"
	"github.com/jessevdk/go-flags"
)

var logger = lw.NewLogger()

func init() {
	parser := flags.NewParser(&conf.Config, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}
}

func neverExit() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			go neverExit()
		}
	}()
	err := shodan.ListenAlerts()
	if err != nil {
		logger.Error(err)
		go neverExit()
	}
}

func main() {
	logger.Debug("main process started")
	go neverExit()
	select {}
}
