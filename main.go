package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"pluto/global"
	"pluto/mapping"
	"pluto/util"
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
	//global.CheckLibrary()
	//Test
	m, err := mapping.LoadMapping("1.20.1", "yarn")
	if err != nil {
		log.Fatal(err)
	}
	count := 0
	for notch, official := range m {
		fmt.Printf("Notch: n=%s c=%s s=%s-> Official: n=%s c=%s s=%s\n", notch.Name, notch.Class, notch.Signature, official.Name, official.Class, official.Signature)
		count++
		if count >= 50 {
			break
		}
	}

	//Main Logic
	//err = webserver.Launch()
	//if err != nil {
	//	log.Fatal(err)
	//}
}
