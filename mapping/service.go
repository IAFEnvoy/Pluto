package mapping

import (
	"errors"
	"pluto/global"
	"pluto/mapping/misc"
	"pluto/mapping/services"
	"pluto/util"
)

type Service interface {
	GetName() string
	GetPathOrDownload(mcVersion string) (string, error)
	LoadMapping(mcVersion string) (map[misc.SingleInfo]misc.SingleInfo, error) //All default is notch->target
	Remap(mcVersion string) (string, error)
}

var serviceMap = map[string]Service{
	"official": &services.Official{},
	"yarn":     &services.Yarn{},
}

func LoadMapping(mcVersion, mapping string) (map[misc.SingleInfo]misc.SingleInfo, error) {
	service, ok := serviceMap[mapping]
	if !ok {
		return nil, errors.New("unknown mapping type")
	}
	return service.LoadMapping(mcVersion)
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
