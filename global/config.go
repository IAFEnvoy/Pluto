package global

import (
	encoder "github.com/zwgblue/yaml-encoder"
	"gopkg.in/yaml.v3"
	"os"
	"pluto/util"
)

type Urls struct {
	MavenCentral       string `yaml:"MavenCentral"`
	FabricMaven        string `yaml:"fabricMaven"`
	FabricMeta         string `yaml:"fabricMeta"`
	MojangLauncherMeta string `yaml:"mojangLauncherMeta"`
}

type ConfigObject struct {
	JavaPath string `yaml:"javaPath" comment:"executable java file for command"`
	Urls     Urls   `yaml:"urls"`
}

const ConfigPath = "config.yaml"

var Config = ConfigObject{
	JavaPath: "java",
	Urls: Urls{
		MavenCentral:       "https://repo1.maven.org/maven2",
		FabricMaven:        "https://maven.fabricmc.net",
		FabricMeta:         "https://meta.fabricmc.net",
		MojangLauncherMeta: "https://launchermeta.mojang.com",
	},
}

func LoadConfig() error {
	util.LOGGER.Info("Loading global config: " + ConfigPath)
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
