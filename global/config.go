package global

import (
	"github.com/gin-contrib/cors"
	encoder "github.com/zwgblue/yaml-encoder"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

type Urls struct {
	MavenCentral       string `yaml:"MavenCentral"`
	MojangLauncherMeta string `yaml:"mojangLauncherMeta"`
	MojangPistonData   string `yaml:"mojangPistonData"`
	MojangPistonMeta   string `yaml:"mojangPistonMeta"`
	FabricMaven        string `yaml:"fabricMaven"`
	FabricMeta         string `yaml:"fabricMeta"`
	NeoForgeMaven      string `yaml:"neoForgeMaven"`
}

type JavaProgramConfig struct {
	JavaParams       []string `yaml:"javaParams"`
	DecompilerParams []string `yaml:"decompilerParams"`
}

type ConfigObject struct {
	Port       int               `yaml:"port" comment:"http server port"`
	JavaPath   string            `yaml:"javaPath" comment:"executable java file for command"`
	Urls       Urls              `yaml:"urls" comment:"if official source is too slow, try BMCLAPI: https://bmclapidoc.bangbang93.com/"`
	Remapper   JavaProgramConfig `yaml:"remapper"`
	Decompiler JavaProgramConfig `yaml:"decompiler"`
	Cors       cors.Config       `yaml:"cors"`
}

const ConfigPath = "config.yaml"

var Config = ConfigObject{
	Port:     5678,
	JavaPath: "java",
	Urls: Urls{
		MavenCentral:       "https://repo1.maven.org/maven2",
		MojangLauncherMeta: "https://launchermeta.mojang.com",
		MojangPistonData:   "https://piston-data.mojang.com",
		MojangPistonMeta:   "https://piston-meta.mojang.com",
		FabricMaven:        "https://maven.fabricmc.net",
		FabricMeta:         "https://meta.fabricmc.net",
		NeoForgeMaven:      "https://maven.neoforged.net/releases",
	},
	Remapper: JavaProgramConfig{
		JavaParams:       []string{"-Xms2G", "-Xmx2G"},
		DecompilerParams: []string{},
	},
	Decompiler: JavaProgramConfig{
		JavaParams:       []string{"-Xms2G", "-Xmx2G"},
		DecompilerParams: []string{"--thread-count=1", "--skip-extra-files"},
	},
	Cors: cors.Config{
		AllowOrigins: []string{"localhost"},
		AllowMethods: []string{"GET"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
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
