package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/coreos/go-systemd/daemon"
	"github.com/macarrie/flemzerd/configuration"
	log "github.com/macarrie/flemzerd/logging"
	flag "github.com/ogier/pflag"

	"github.com/macarrie/flemzerd/scheduler"
	"github.com/macarrie/flemzerd/server"
)

func main() {
	debugMode := flag.BoolP("debug", "d", false, "Start in debug mode")
	versionFlag := flag.BoolP("version", "v", false, "Display version number")
	configFilePath := flag.StringP("config", "c", "", "Configuration file path to use")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("flemzerd version %s", configuration.Version)
		os.Exit(0)
	}

	if *debugMode {
		log.Setup(true)
	} else {
		log.Setup(false)
	}

	log.Info("Starting flemzerd (version ", configuration.Version, ")")

	if *configFilePath != "" {
		log.Info("Loading provided configuration file")
		configuration.UseFile(*configFilePath)
	}

	scheduler.Run(*debugMode)
	server.Stop()
	if configuration.Config.Interface.Enabled {
		go server.Start(configuration.Config.Interface.Port)
	}
	daemon.SdNotify(false, "READY=1")

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	for {
		switch sig := <-signalChannel; sig {
		case syscall.SIGINT, syscall.SIGTERM:
			log.Info("Shutting down...")
			server.Stop()
			scheduler.Stop()
			os.Exit(0)
		case syscall.SIGUSR1:
			log.Info("Signal received: reloading configuration")
			daemon.SdNotify(false, "READY=0")

			server.Stop()
			scheduler.Reload(*debugMode)
			if configuration.Config.Interface.Enabled {
				go server.Start(configuration.Config.Interface.Port)
			}

			daemon.SdNotify(false, "READY=1")
		}
	}
}
