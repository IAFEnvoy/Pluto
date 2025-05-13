package services

import (
	"os"
	"pluto/global"
	"pluto/util"
	"pluto/util/network"
	"pluto/vanilla"
)

type Official struct{}

func (s *Official) GetName() string {
	return "official"
}

func (s *Official) GetPathOrDownload(mcVersion string) (string, error) {
	path := global.GetMappingPath(s, mcVersion, "txt")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return path, nil
	}
	downloads, err := vanilla.GetOrDownload(mcVersion)
	if err != nil {
		return "", err
	}
	data, err := network.Get(downloads.ClientMappings.Url)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(path, data, 0666)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (s *Official) Remap(mcVersion string) (string, error) {
	jarPath, err := vanilla.GetMcJarPath(mcVersion)
	if err != nil {
		return "", err
	}
	mappingPath, err := s.GetPathOrDownload(mcVersion)
	if err != nil {
		return "", err
	}
	outputPath := global.GetRemappedPath(s, mcVersion)
	util.ExecuteCommand(global.Config.JavaPath, []string{"-cp", global.ClassPath, global.ArtMainClass, "--input", jarPath, "--output", outputPath, "--map", mappingPath, "--reverse"}, true)
	return outputPath, nil
}
