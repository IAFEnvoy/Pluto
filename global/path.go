package global

import (
	"os"
	"path/filepath"
)

type Named interface {
	GetName() string
}

type NamedImpl struct {
	Name string
}

func (n NamedImpl) GetName() string {
	return n.Name
}

func CreatePathAndReturn(path, file string) string {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return filepath.Join(path, file)
}

func GetMappingPath(named Named, mcVersion, extension string) string {
	return CreatePathAndReturn(filepath.Join("cache", "mappings", named.GetName()), mcVersion+"."+extension)
}

func GetMinecraftPath(mcVersion string) string {
	return CreatePathAndReturn(filepath.Join("cache", "minecraft"), mcVersion+".jar")
}

func GetRemappedPath(named Named, mcVersion string) string {
	return CreatePathAndReturn(filepath.Join("cache", "remapped", named.GetName()), mcVersion+".jar")
}

func GetSourceFolder(named Named, mcVersion string) string {
	path := filepath.Join("cache", "remapped", named.GetName(), mcVersion)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return path
}
