package download

import (
	"os"
	"pluto/util/network"
)

const officialMappingFolder = "cache/mappings/official/"

func GetOfficial(mcVersion string) ([]byte, error) {
	err := os.MkdirAll(officialMappingFolder, os.ModePerm)
	if err != nil {
		return nil, err
	}
	path := officialMappingFolder + mcVersion + ".txt"
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	downloads, err := network.GetOrDownload(mcVersion)
	if err != nil {
		return nil, err
	}
	data, err := network.Get(downloads.ClientMappings.Url)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(path, data, 0666)
	if err != nil {
		return nil, err
	}
	return data, nil
}
