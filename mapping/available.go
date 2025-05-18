package mapping

import (
	"pluto/util"
)

type AvailableConfig struct {
	Official []string `json:"official"`
	Yarn     []string `json:"yarn"`
}

type TaskInfo struct {
	MappingType string
	Version     string
}

const configPath = "cache/available.json"

var (
	availableConfig = AvailableConfig{}
	pendingTasks    []TaskInfo
)

func InitMappingConfig() error {
	config, err := util.LoadConfig[AvailableConfig](configPath)
	if err != nil {
		return err
	}
	availableConfig = config
	return nil
}

func IsAvailable(mcVersion, mappingType string) bool {
	switch mappingType {
	case "official":
		return util.Contains(availableConfig.Official, mcVersion)
	case "yarn":
		return util.Contains(availableConfig.Yarn, mcVersion)
	default:
		return false
	}
}

func IsPending(mcVersion, mappingType string) bool {
	return util.Contains(pendingTasks, TaskInfo{mappingType, mcVersion})
}

func CanAddTask(mcVersion, mappingType string) bool {
	return !IsAvailable(mcVersion, mappingType) && !IsPending(mcVersion, mappingType)
}
