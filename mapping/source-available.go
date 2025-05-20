package mapping

import (
	"log/slog"
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

const configPath = "cache/source-available.json"

var (
	availableConfig = AvailableConfig{}
	pendingTasks    = make(map[TaskInfo]struct{})
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

func FailurePending(mcVersion, mappingType string) {
	delete(pendingTasks, TaskInfo{
		MappingType: mappingType,
		Version:     mcVersion,
	})
}

func Done(mcVersion, mappingType string) {
	FailurePending(mcVersion, mappingType)
	switch mappingType {
	case "official":
		availableConfig.Official = append(availableConfig.Official, mcVersion)
	case "yarn":
		availableConfig.Yarn = append(availableConfig.Yarn, mcVersion)
	}
	err := util.SaveConfig(availableConfig, configPath)
	if err != nil {
		slog.Error("Failed to save " + configPath + ": " + err.Error())
	}
}
