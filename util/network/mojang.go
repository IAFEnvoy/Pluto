package network

import (
	"encoding/json"
	"errors"
	"pluto/global"
)

type SingleManifest struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type VersionManifest struct {
	Versions []SingleManifest `json:"versions"`
}

type SingleFile struct {
	Sha1 string `json:"sha1"`
	Size int64  `json:"size"`
	Url  string `json:"url"`
}

type Downloads struct {
	Client         SingleFile     `json:"client"`
	ClientMappings SingleManifest `json:"client_mappings"`
	Server         SingleFile     `json:"server"`
	ServerMappings SingleManifest `json:"server_mappings"`
}

type PistonData struct {
	Downloads Downloads `json:"downloads"`
}

var cache = map[string]Downloads{}

func GetOrDownload(mcVersion string) (Downloads, error) {
	if downloads, ok := cache[mcVersion]; ok {
		return downloads, nil
	}
	//request launcher meta
	data, err := Get(global.Config.Urls.MojangLauncherMeta + "/mc/game/version_manifest_v2.json")
	if err != nil {
		return Downloads{}, err
	}
	manifest := VersionManifest{}
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return Downloads{}, err
	}
	var url = ""
	for _, version := range manifest.Versions {
		if version.Id == mcVersion {
			url = version.Url
			break
		}
	}
	if url == "" {
		return Downloads{}, errors.New("Cannot find mc version " + mcVersion)
	}
	//request piston data
	data, err = Get(url)
	if err != nil {
		return Downloads{}, err
	}
	downloads := PistonData{}
	err = json.Unmarshal(data, &downloads)
	if err != nil {
		return Downloads{}, err
	}
	cache[mcVersion] = downloads.Downloads
	return downloads.Downloads, nil
}
