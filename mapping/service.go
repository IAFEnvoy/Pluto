package mapping

import (
	"errors"
	"fmt"
	"log/slog"
	"pluto/global"
	"pluto/mapping/java"
	"pluto/mapping/services"
	"pluto/util"
)

type Service interface {
	GetName() string
	GetPathOrDownload(mcVersion string) (string, error)
	GetMappingCacheOrError(mcVersion string) (*java.Mappings, error)
	SaveMappingCache(mcVersion string, mapping *java.Mappings)
	LoadMapping(mcVersion string) (*map[java.SingleInfo]java.SingleInfo, error) //All default is notch->target
	Remap(mcVersion string) (string, error)
}

var (
	serviceMap = map[string]Service{
		"official": &services.Official{},
		"yarn":     &services.Yarn{},
	}
	loadMappingLock = util.NewNamedLock()
)

func CachedMapping(mcVersion, mappingType string) bool {
	service, ok := serviceMap[mappingType]
	if !ok {
		return false
	}
	_, err := service.GetMappingCacheOrError(mcVersion)
	return err == nil
}

func LoadMapping(mcVersion, mappingType string) (*java.Mappings, error) {
	service, ok := serviceMap[mappingType]
	if !ok {
		return &java.Mappings{}, errors.New("unknown mapping type")
	}
	if cache, err := service.GetMappingCacheOrError(mcVersion); err == nil {
		return cache, nil
	}

	if loadMappingLock.IsLocked(mcVersion, mappingType) {
		slog.Warn("This mapping is loading!")
		return nil, errors.New("this mapping is loading")
	}
	loadMappingLock.Lock(mcVersion, mappingType)
	defer loadMappingLock.Unlock(mcVersion, mappingType)

	slog.Info(fmt.Sprintf("Loading mapping type %s for %s", mappingType, mcVersion))
	m, err := service.LoadMapping(mcVersion)
	if err != nil {
		return &java.Mappings{}, err
	}
	m3 := java.BuildMapping(m)
	service.SaveMappingCache(mcVersion, m3)
	return m3, nil
}

func GenerateSource(mcVersion, mappingType string) (string, error) {
	if !CanAddTask(mcVersion, mappingType) {
		return "", errors.New("this type has generated or generating")
	}
	service, ok := serviceMap[mappingType]
	if !ok {
		return "", errors.New("unknown mapping type")
	}
	StartPending(mcVersion, mappingType)
	path, err := service.Remap(mcVersion)
	if err != nil {
		return "", err
	}
	sourcePath := global.GetSourceFolder(service, mcVersion)
	err = util.ExecuteCommand(global.Config.JavaPath, []string{"-jar", global.DecompilerPath, path, sourcePath}, true)
	if err != nil {
		return "", err
	}
	Done(mcVersion, mappingType)
	return sourcePath, nil
}
