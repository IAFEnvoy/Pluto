package global

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"pluto/util"
	"strconv"
	"strings"
	"time"
)

type ToolInfo struct {
	Name            string `json:"name"`
	CurrentVersion  string `json:"currentVersion"`
	LastChecked     string `json:"lastChecked"`
	Hash            string `json:"hash"`
	URL             string `json:"url"`
	MavenGroupID    string `json:"mavenGroupID,omitempty"`
	MavenArtifactID string `json:"mavenArtifactID,omitempty"`
	MavenRepoURL    string `json:"mavenRepoURL,omitempty"`
}

type LibraryConfig struct {
	Tools []ToolInfo `json:"libraries"`
}

const (
	LibraryPath           = "libraries"
	configPath            = LibraryPath + "/versions.json"
	httpTimeout           = 30 * time.Second
	TinyRemapperMainClass = "net.fabricmc.tinyremapper.Main"
)

var (
	ClassPath      = LibraryPath + "/*"
	DecompilerPath = filepath.Join(LibraryPath, "vineflower.jar")
)

func CheckLibrary() {
	util.LOGGER.Info("Start checking libraries")
	// 确保保存工具的目录存在
	if err := os.MkdirAll(LibraryPath, 0755); err != nil {
		util.LOGGER.Error("Failed to create library folder:  " + err.Error())
		return
	}

	// 加载配置文件
	config, err := loadConfig()
	if err != nil {
		util.LOGGER.Error("Failed to load library global:  " + err.Error())
		config = LibraryConfig{}
	}

	// 检查并更新工具
	toolsToUpdate := []ToolInfo{
		{
			Name:            "vineflower.jar",
			MavenGroupID:    "org.vineflower",
			MavenArtifactID: "vineflower",
			MavenRepoURL:    Config.Urls.MavenCentral,
		},
		{
			Name:            "tiny-remapper.jar",
			MavenGroupID:    "net.fabricmc",
			MavenArtifactID: "tiny-remapper",
			MavenRepoURL:    Config.Urls.FabricMaven,
		},
		{
			Name:            "mapping-io.jar",
			MavenGroupID:    "net.fabricmc",
			MavenArtifactID: "mapping-io",
			MavenRepoURL:    Config.Urls.FabricMaven,
		},
		{
			Name:            "asm.jar",
			MavenGroupID:    "org.ow2.asm",
			MavenArtifactID: "asm",
			MavenRepoURL:    Config.Urls.FabricMaven,
			CurrentVersion:  "9.8",
		},
		{
			Name:            "asm-commons.jar",
			MavenGroupID:    "org.ow2.asm",
			MavenArtifactID: "asm-commons",
			MavenRepoURL:    Config.Urls.FabricMaven,
			CurrentVersion:  "9.8",
		},
		{
			Name:            "asm-tree.jar",
			MavenGroupID:    "org.ow2.asm",
			MavenArtifactID: "asm-tree",
			MavenRepoURL:    Config.Urls.FabricMaven,
			CurrentVersion:  "9.8",
		},
	}

	for _, tool := range toolsToUpdate {
		if err := checkAndUpdateTool(tool, &config); err != nil {
			util.LOGGER.Error("Failed to update " + tool.Name + ":  " + err.Error())
		}
	}

	// 保存更新后的配置
	if err := saveConfig(config); err != nil {
		util.LOGGER.Error("Failed to save library global:  " + err.Error())
	}
}

// 加载配置文件
func loadConfig() (LibraryConfig, error) {
	var config LibraryConfig

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	// 解析配置文件
	if err := json.Unmarshal(data, &config); err != nil {
		return config, err
	}

	return config, nil
}

// 保存配置文件
func saveConfig(config LibraryConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// 检查并更新工具
func checkAndUpdateTool(tool ToolInfo, config *LibraryConfig) error {
	// 查找工具在配置中的索引
	index := -1
	for i, t := range config.Tools {
		if t.Name == tool.Name {
			index = i
			break
		}
	}

	// 获取当前时间
	now := time.Now().Format(time.RFC3339)

	// 获取最新版本信息
	var latestVersion string
	var err error

	if tool.MavenGroupID != "" && tool.MavenArtifactID != "" {
		// 从Maven获取最新版本
		latestVersion, err = getLatestMavenVersion(tool.MavenGroupID, tool.MavenArtifactID, tool.MavenRepoURL)
		if err != nil {
			latestVersion = tool.CurrentVersion
			util.LOGGER.Error("Failed to load versions, try to use hardcoded version " + latestVersion)
		}
		tool.URL = fmt.Sprintf("%s/%s/%s/%s/%s-%s.jar", tool.MavenRepoURL, strings.ReplaceAll(tool.MavenGroupID, ".", "/"), tool.MavenArtifactID, latestVersion, tool.MavenArtifactID, latestVersion)
	} else if tool.URL != "" {
		latestVersion, err = getVersionFromURL(tool.URL)
		if err != nil {
			util.LOGGER.Error("Failed to load versions:  " + err.Error())
			latestVersion = tool.CurrentVersion
		}
	} else {
		util.LOGGER.Error("Cannot update " + tool.Name + " since no url provided")
		return err
	}

	// 检查是否需要更新
	needsUpdate := false
	currentVersion := ""
	toolPath := filepath.Join(LibraryPath, tool.Name)

	// 检查本地是否已有该工具
	if index >= 0 {
		currentVersion = config.Tools[index].CurrentVersion

		// 比较版本
		if currentVersion != latestVersion {
			needsUpdate = true
		}
	} else {
		// 本地没有该工具，需要下载
		needsUpdate = true
	}

	// 如果需要更新，下载最新版本
	if needsUpdate {
		util.LOGGER.Info("Updating " + tool.Name + " to " + latestVersion + "...")

		// 下载文件
		tmpPath := toolPath + ".tmp"
		if err := downloadFile(tmpPath, tool.URL); err != nil {
			util.LOGGER.Error("Failed to download file:  " + err.Error())
			return err
		}

		// 计算下载文件的哈希值
		newHash, err := calculateFileHash(tmpPath)
		if err != nil {
			util.LOGGER.Error("Failed to calculate file hash:  " + err.Error())
			return err
		}

		// 替换旧文件
		if err := os.Rename(tmpPath, toolPath); err != nil {
			util.LOGGER.Error("Failed to replace file:  " + err.Error())
			return err
		}

		util.LOGGER.Info("Successfully update " + tool.Name + " to " + latestVersion)

		// 更新配置
		newTool := ToolInfo{
			Name:            tool.Name,
			CurrentVersion:  latestVersion,
			LastChecked:     now,
			Hash:            newHash,
			URL:             tool.URL,
			MavenGroupID:    tool.MavenGroupID,
			MavenArtifactID: tool.MavenArtifactID,
			MavenRepoURL:    tool.MavenRepoURL,
		}

		if index >= 0 {
			config.Tools[index] = newTool
		} else {
			config.Tools = append(config.Tools, newTool)
		}
	} else {
		// 文件是最新的，只更新检查时间
		if index >= 0 {
			config.Tools[index].LastChecked = now
		} else {
			// 添加工具信息到配置
			config.Tools = append(config.Tools, ToolInfo{
				Name:            tool.Name,
				CurrentVersion:  currentVersion,
				LastChecked:     now,
				Hash:            "", // 没有计算哈希
				URL:             tool.URL,
				MavenGroupID:    tool.MavenGroupID,
				MavenArtifactID: tool.MavenArtifactID,
				MavenRepoURL:    tool.MavenRepoURL,
			})
		}

		util.LOGGER.Info(tool.Name + " is already the latest version " + currentVersion)
	}

	return nil
}

// 从Maven仓库获取最新版本
func getLatestMavenVersion(groupID, artifactID, repoURL string) (string, error) {
	// 将 groupID 中的点转换为路径分隔符
	groupPath := strings.ReplaceAll(groupID, ".", "/")

	// 构建 metadata.xml 的URL
	metadataURL := fmt.Sprintf("%s/%s/%s/maven-metadata.xml", repoURL, groupPath, artifactID)

	// 获取 metadata.xml 内容
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Get(metadataURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		util.LOGGER.Error("Failed to load versions from " + metadataURL + " with status code " + strconv.Itoa(resp.StatusCode))
		return "", errors.New("failed to load versions from maven")
	}

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析版本号
	metadata := string(body)

	// 尝试提取<latest>标签
	versionStart := strings.Index(metadata, "<latest>")
	if versionStart != -1 {
		versionStart += 8
		versionEnd := strings.Index(metadata[versionStart:], "</latest>")
		if versionEnd != -1 {
			latestVersion := metadata[versionStart : versionStart+versionEnd]
			return latestVersion, nil
		}
	}

	// 如果没有<latest>标签，尝试提取<version>标签中的最后一个版本
	versionStart = strings.Index(metadata, "<version>")
	if versionStart == -1 {
		util.LOGGER.Error("Failed to parse latest version")
		return "", errors.New("failed to parse latest version")
	}

	var versions []string
	for versionStart != -1 {
		versionStart += 9
		versionEnd := strings.Index(metadata[versionStart:], "</version>")
		if versionEnd == -1 {
			break
		}

		versions = append(versions, metadata[versionStart:versionStart+versionEnd])
		metadata = metadata[versionStart+versionEnd+10:]
		versionStart = strings.Index(metadata, "<version>")
	}

	if len(versions) == 0 {
		util.LOGGER.Error("No versions found for " + artifactID)
		return "", errors.New("no versions found for " + artifactID)
	}
	latestVersion := versions[len(versions)-1]
	return latestVersion, nil
}

// 从URL获取版本信息
func getVersionFromURL(url string) (string, error) {
	// 对于 FernFlower，我们无法直接获取版本，所以使用时间戳作为版本
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Head(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		util.LOGGER.Error("Failed to load url head from maven with status code " + strconv.Itoa(resp.StatusCode))
		return "", errors.New("failed to load url head from maven")
	}

	// 使用Last-Modified头作为版本标识
	lastModified := resp.Header.Get("Last-Modified")
	if lastModified == "" {
		// 如果没有Last-Modified头，使用当前时间
		return time.Now().Format("20060102150405"), nil
	}

	// 解析Last-Modified时间
	t, err := time.Parse(time.RFC1123, lastModified)
	if err != nil {
		// 解析失败，使用原始字符串
		return lastModified, nil
	}

	// 使用时间戳作为版本
	return t.Format("20060102150405"), nil
}

// 计算文件的 SHA-256 哈希值
func calculateFileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// 下载文件
func downloadFile(dstPath, url string) error {
	util.LOGGER.Info("Downloading: " + url)
	// 创建临时文件
	tmpFile, err := os.CreateTemp(filepath.Dir(dstPath), filepath.Base(dstPath)+".tmp.*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	defer func() {
		err = tmpFile.Close()
		// 如果下载失败，删除临时文件
		if err != nil {
			err := os.Remove(tmpPath)
			if err != nil {
				return
			}
		}
	}()

	// 获取URL内容
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		util.LOGGER.Error("Failed to download from maven with status code " + strconv.Itoa(resp.StatusCode))
		return errors.New("failed to download from maven")
	}

	// 写入文件
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return err
	}

	// 确保文件内容已写入磁盘
	if err := tmpFile.Sync(); err != nil {
		return err
	}

	// 关闭文件
	if err := tmpFile.Close(); err != nil {
		return err
	}

	// 移动临时文件到目标位置
	if err := os.Rename(tmpPath, dstPath); err != nil {
		return err
	}

	return nil
}
