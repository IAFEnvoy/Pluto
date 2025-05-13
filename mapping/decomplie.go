package mapping

import (
	"os"
	"path/filepath"
	"pluto/global"
	"pluto/mapping/download"
	"pluto/util"
)

const RemappedFolder = "cache/remapped/"
const SourceFolder = "cache/source/"

func DecompileSync(mcVersion string, mapping string) {
	util.Execute(func() error {
		jarPath, err := download.GetMcJarPath(mcVersion)
		if err != nil {
			return err
		}
		mappingPath, err := download.GetYarnPath(mcVersion)
		if err != nil {
			return err
		}
		remappedPath := filepath.Join(RemappedFolder, mapping)
		remappedFile := filepath.Join(remappedPath, mcVersion+".jar")
		sourcePath := filepath.Join(SourceFolder, mapping, mcVersion)
		err = os.MkdirAll(remappedPath, os.ModePerm)
		if err != nil {
			return err
		}
		err = os.MkdirAll(sourcePath, os.ModePerm)
		if err != nil {
			return err
		}
		util.ExecuteCommand(global.Config.JavaPath, []string{"-cp", global.ClassPath, global.TinyRemapperMainClass, jarPath, remappedFile, mappingPath, "official", "named"}, false)
		util.LOGGER.Info("Start decompiling, please wait...")
		util.ExecuteCommand(global.Config.JavaPath, []string{"-jar", global.DecompilerPath, remappedFile, sourcePath}, true)
		return nil
	})
}
