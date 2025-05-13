package main

import (
	"log"
	"os"
	"pluto/global"
	"pluto/mapping"
	"pluto/util"
)

func main() {
	defer util.CloseWorkers()
	util.LOGGER.Info("Launching Pluto v" + global.VERSION)
	err := global.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll("temp", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	global.CheckLibrary()

	_, err = mapping.GenerateSource("1.20.1", "official")
	if err != nil {
		util.LOGGER.Error(err.Error())
	}
}
