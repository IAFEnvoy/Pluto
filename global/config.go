package global

import (
	encoder "github.com/zwgblue/yaml-encoder"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

type Urls struct {
	MavenCentral       string `yaml:"MavenCentral"`
	MojangLauncherMeta string `yaml:"mojangLauncherMeta"`
	FabricMaven        string `yaml:"fabricMaven"`
	FabricMeta         string `yaml:"fabricMeta"`
	NeoForgeMaven      string `yaml:"neoForgeMaven"`
}

type ConfigObject struct {
	Port     int    `yaml:"port" comment:"http server port"`
	JavaPath string `yaml:"javaPath" comment:"executable java file for command"`
	Urls     Urls   `yaml:"urls" comment:"if official source is too slow, try BMCLAPI: https://bmclapidoc.bangbang93.com/"`
}

const ConfigPath = "config.yaml"

var Config = ConfigObject{
	Port:     5678,
	JavaPath: "java",
	Urls: Urls{
		MavenCentral:       "https://repo1.maven.org/maven2",
		MojangLauncherMeta: "https://launchermeta.mojang.com",
		FabricMaven:        "https://maven.fabricmc.net",
		FabricMeta:         "https://meta.fabricmc.net",
		NeoForgeMaven:      "https://maven.neoforged.net/releases",
	},
}

func LoadConfig() error {
	slog.Info("Loading global config: " + ConfigPath)
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		encoded := encoder.NewEncoder(Config, encoder.WithComments(encoder.CommentsOnHead))
		bytes, err := encoded.Encode()
		if err != nil {
			return err
		}
		err = os.WriteFile(ConfigPath, bytes, 0644)
		if err != nil {
			return err
		}
	} else {
		data, err := os.ReadFile(ConfigPath)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(data, &Config)
		if err != nil {
			return err
		}
	}
	return nil
}
