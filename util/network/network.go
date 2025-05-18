package network

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"pluto/util"
	"strconv"
	"time"
)

func Get(url string) ([]byte, error) {
	slog.Info("Downloading: " + url)
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.Itoa(resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取 Yarn 版本列表响应失败: %v", err)
	}
	return body, nil
}

func File(url string, path string) error {
	slog.Info("Downloading: " + url)
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(strconv.Itoa(resp.StatusCode))
	}

	// 创建临时文件
	tempFile, err := os.CreateTemp("temp", "download-*.tmp")
	if err != nil {
		return err
	}
	tempPath := tempFile.Name()
	defer func() {
		tempFile.Close()
		os.Remove(tempPath)
	}()

	file, err := os.Create(tempPath)
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	// 确保所有数据都写入磁盘
	err = tempFile.Sync()
	if err != nil {
		return err
	}
	tempFile.Close()
	_, err = util.CopyFile(tempFile.Name(), path)
	if err != nil {
		return err
	}
	return err
}
