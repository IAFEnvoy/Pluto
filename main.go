package main

import (
	"log"
	"os"
	"pluto/global"
	"pluto/mapping"
	"pluto/mapping/download"
	"pluto/util"
)

func main() {
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
	_, err = download.GetYarnPath("1.20.1")
	if err != nil {
		util.LOGGER.Error("Unable to download yarn: " + err.Error())
		return
	}
	_, err = download.GetMcJarPath("1.20.1")
	if err != nil {
		util.LOGGER.Error("Unable to download mcjar: " + err.Error())
		return
	}
	_, err = download.GetOfficial("1.20.1")
	if err != nil {
		util.LOGGER.Error("Unable to download official: " + err.Error())
		return
	}
	mapping.DecompileSync("1.20.1", "yarn")
	util.CloseWorkers()
}
