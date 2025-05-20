package main

import (
	"log"
	"log/slog"
	"os"
	"pluto/global"
	"pluto/mapping"
	"pluto/util"
	"pluto/webserver"
)

func main() {
	defer util.CloseWorkers()
	util.InitLogger()
	slog.Info("Launching Pluto v" + global.Version)
	//Config
	slog.Info("Loading configs...")
	err := global.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = mapping.InitMappingConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll("temp", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	//Libraries
	global.CheckLibrary()
	//Main Logic
	err = webserver.Launch()
	if err != nil {
		log.Fatal(err)
	}
}
