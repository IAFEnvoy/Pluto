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
	pendingTasks    map[TaskInfo]struct{}
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
	_, ok := pendingTasks[TaskInfo{
		MappingType: mappingType,
		Version:     mcVersion,
	}]
	return ok
}

func CanAddTask(mcVersion, mappingType string) bool {
	return !IsAvailable(mcVersion, mappingType) && !IsPending(mcVersion, mappingType)
}

func StartPending(mcVersion, mappingType string) {
	pendingTasks[TaskInfo{
		MappingType: mappingType,
		Version:     mcVersion,
	}] = struct{}{}
}

func Done(mcVersion, mappingType string) {
	delete(pendingTasks, TaskInfo{
		MappingType: mappingType,
		Version:     mcVersion,
	})
	switch mappingType {
	case "official":
		availableConfig.Official = append(availableConfig.Official, mcVersion)
	case "yarn":
		availableConfig.Yarn = append(availableConfig.Yarn, mcVersion)
	}
}
