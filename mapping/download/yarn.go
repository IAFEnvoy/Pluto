package download

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"pluto/global"
	"pluto/util"
	"pluto/util/network"
)

type YarnVersion struct {
	GameVersion string `json:"gameVersion"`
	Separator   string `json:"separator"`
	Build       int    `json:"build"`
	Maven       string `json:"maven"`
	Version     string `json:"version"`
	Stable      bool   `json:"stable"`
}

const yarnMappingFolder = "cache/mappings/yarn/"

func GetYarnPath(mcVersion string) (string, error) {
	err := os.MkdirAll(yarnMappingFolder, os.ModePerm)
	if err != nil {
		return "", err
	}
	path := yarnMappingFolder + mcVersion + ".tiny"
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return path, nil
	}
	body, err := network.Get(global.Config.Urls.FabricMeta + "/v2/versions/yarn")
	if err != nil {
		util.LOGGER.Error("Unable to download yarn version: " + err.Error())
		return "", err
	}
	var versions []YarnVersion
	if err := json.Unmarshal(body, &versions); err != nil {
		util.LOGGER.Error("Unable to unmarshal yarn versions: " + err.Error())
		return "", err
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
		util.LOGGER.Error("Unable to find latest version for " + mcVersion)
		return "", fmt.Errorf("unable to find latest version for %s", mcVersion)
	}
	jar, err := network.Get(fmt.Sprintf(global.Config.Urls.FabricMaven+"/net/fabricmc/yarn/%s/yarn-%s-tiny.gz", latestVersion.Version, latestVersion.Version))
	if err != nil {
		util.LOGGER.Error("Unable to download yarn mapping: " + err.Error())
		return "", err
	}
	data, err := getMappingsTinyFromGzip(jar)
	if err != nil {
		util.LOGGER.Error("Unable to unzip yarn mapping: " + err.Error())
		return "", err
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return "", err
	}
	return path, nil
}

func getMappingsTinyFromGzip(gzipData []byte) ([]byte, error) {
	gzipReader, err := gzip.NewReader(bytes.NewReader(gzipData))
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()
	content, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, err
	}
	return content, nil
}
