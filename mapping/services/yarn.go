package services

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"pluto/global"
	"pluto/util"
	"pluto/util/network"
	"pluto/vanilla"
)

type Yarn struct{}

type YarnVersion struct {
	GameVersion string `json:"gameVersion"`
	Separator   string `json:"separator"`
	Build       int    `json:"build"`
	Maven       string `json:"maven"`
	Version     string `json:"version"`
	Stable      bool   `json:"stable"`
}

func (s Yarn) GetName() string {
	return "yarn"
}

func (s Yarn) GetPathOrDownload(mcVersion string) (string, error) {
	path := global.GetMappingPath(s, mcVersion, "tiny")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return path, nil
	}
	body, err := network.Get(global.Config.Urls.FabricMeta + "/v2/versions/yarn")
	if err != nil {
		return "", errors.New("Unable to download yarn versions: " + err.Error())
	}
	var versions []YarnVersion
	if err := json.Unmarshal(body, &versions); err != nil {
		return "", errors.New("Unable to unmarshal yarn versions: " + err.Error())
	}
	var latestVersion *YarnVersion
	for i := range versions {
		version := &versions[i]
		if version.GameVersion == mcVersion {
			if latestVersion == nil || version.Build > latestVersion.Build {
				latestVersion = version
			}
		}
	}
	if latestVersion == nil {
		return "", errors.New("Unable to find latest version for " + mcVersion)
	}
	jar, err := network.Get(fmt.Sprintf(global.Config.Urls.FabricMaven+"/net/fabricmc/yarn/%s/yarn-%s-tiny.gz", latestVersion.Version, latestVersion.Version))
	if err != nil {
		return "", errors.New("Unable to download yarn mapping: " + err.Error())
	}
	data, err := getMappingsTinyFromGzip(jar)
	if err != nil {
		return "", errors.New("Unable to unzip yarn mapping: " + err.Error())
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (s Yarn) Remap(mcVersion string) (string, error) {
	jarPath, err := vanilla.GetMcJarPath(mcVersion)
	if err != nil {
		return "", err
	}
	mappingPath, err := s.GetPathOrDownload(mcVersion)
	if err != nil {
		return "", err
	}
	outputPath := global.GetRemappedPath(s, mappingPath)
	util.ExecuteCommand(global.Config.JavaPath, []string{"-cp", global.ClassPath, global.TinyRemapperMainClass, jarPath, outputPath, mappingPath, "official", "named"}, false)
	return outputPath, nil
}

func getMappingsTinyFromGzip(gzipData []byte) ([]byte, error) {
	gzipReader, err := gzip.NewReader(bytes.NewReader(gzipData))
	if err != nil {
		return nil, err
	}
	defer func(gzipReader *gzip.Reader) {
		err := gzipReader.Close()
		if err != nil {
			util.LOGGER.Error("Error closing gzip reader: " + err.Error())
		}
	}(gzipReader)
	content, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, err
	}
	return content, nil
}
