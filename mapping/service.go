package mapping

import (
	"errors"
	"pluto/global"
	"pluto/mapping/services"
	"pluto/util"
)

type Service interface {
	GetName() string
	GetPathOrDownload(mcVersion string) (string, error)
	Remap(mcVersion string) (string, error)
}

var serviceMap = map[string]Service{
	"official": &services.Official{},
	"yarn":     &services.Yarn{},
}

func GenerateSource(mcVersion, mapping string) (string, error) {
	service, ok := serviceMap[mapping]
	if !ok {
		return "", errors.New("unknown mapping type")
	}
	path, err := service.Remap(mcVersion)
	if err != nil {
		return "", err
	}
	sourcePath := global.GetSourceFolder(service, mcVersion)
	util.ExecuteCommand(global.Config.JavaPath, []string{"-jar", global.DecompilerPath, path, sourcePath}, true)
	return sourcePath, nil
}
