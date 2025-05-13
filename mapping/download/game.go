package download

import (
	"os"
	"pluto/util"
	"pluto/util/network"
)

const JarFolder = "cache/minecraft/"

func GetMcJarPath(mcVersion string) (string, error) {
	err := os.MkdirAll(JarFolder, os.ModePerm)
	if err != nil {
		return "", err
	}
	path := JarFolder + mcVersion + ".jar"
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return path, nil
	}
	downloads, err := network.GetOrDownload(mcVersion)
	if err != nil {
		util.LOGGER.Error("Unable to download " + mcVersion + " meta : " + err.Error())
		return "", err
	}
	err = network.File(downloads.Client.Url, path)
	if err != nil {
		util.LOGGER.Error("Unable to download " + mcVersion + " file : " + err.Error())
		return "", err
	}
	return path, nil
}
