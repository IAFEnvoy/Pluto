package services

import (
	"bufio"
	"fmt"
	"os"
	"pluto/convert"
	"pluto/global"
	"pluto/mapping/misc"
	"pluto/util"
	"pluto/util/network"
	"pluto/vanilla"
	"strings"
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

func (s *Official) LoadMapping(mcVersion string) (map[misc.SingleInfo]misc.SingleInfo, error) {
	path, err := s.GetPathOrDownload(mcVersion)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	mapping := make(map[misc.SingleInfo]misc.SingleInfo)
	scanner := bufio.NewScanner(file)
	cachedNotchClass, cachedNamedClass := misc.SingleInfo{}, misc.SingleInfo{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}
		split := strings.Split(line, " -> ")
		if len(split) != 2 {
			continue
		}
		if !strings.HasPrefix(split[0], "    ") { //class
			notch, named := misc.PackClassInfo(strings.ReplaceAll(split[1], ":", "")), misc.PackClassInfo(split[0])
			mapping[notch] = named
			cachedNotchClass, cachedNamedClass = notch, named
		} else if strings.Contains(split[0], ":") { //method
			s := strings.Split(split[0], ":")
			if len(s) == 3 {
				name := strings.Split(strings.Split(s[2], " ")[1], "(")[0]
				signature, err := convert.MethodToByteCodeSignature(s[2], false)
				if err != nil {
					util.Logger.Error("convert to java signature error: " + err.Error())
					continue
				}
				notch, named := misc.PackMethodInfo(split[1], cachedNotchClass.Signature, ""), misc.PackMethodInfo(name, cachedNamedClass.Signature, signature)
				mapping[notch] = named
			}
		} else { //Field
			s := strings.Split(strings.TrimSpace(split[0]), " ")
			if len(s) == 2 {
				notch, named := misc.PackFieldInfo(split[1], cachedNotchClass.Signature, ""), misc.PackFieldInfo(s[1], cachedNamedClass.Signature, s[0])
				mapping[notch] = named
			}
		}
	}
	//End processor
	result, classMapping := make(map[misc.SingleInfo]misc.SingleInfo), make(map[string]string)
	for notch, named := range mapping {
		if notch.Type == "class" {
			classMapping[named.Signature] = notch.Signature
		}
	}
	for notch, named := range mapping {
		notch.Name = convert.FullToClassName(notch.Name)
		named.Name = convert.FullToClassName(named.Name)
		if notch.Type == "method" {
			notch.Signature = misc.ObfuscateMethodSignature(named.Signature, classMapping)
		}
		result[notch] = named
	}
	return result, nil
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
