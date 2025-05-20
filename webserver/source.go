package webserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"pluto/global"
	"pluto/mapping"
	"pluto/util"
)

func initSourceApi(g *gin.Engine) {
	g.GET("/api/source/load", func(c *gin.Context) {
		mcVersion, mappingType := c.Query("version"), c.Query("type")
		if mcVersion == "" || mappingType == "" {
			c.String(http.StatusBadRequest, "Missing query parameter(s)")
			return
		}
		if mapping.IsAvailable(mcVersion, mappingType) {
			c.String(http.StatusOK, "Decompiled")
			return
		}
		if mapping.IsPending(mcVersion, mappingType) {
			c.String(http.StatusForbidden, "This task is pending")
			return
		}
		util.Execute(func() error {
			_, err := mapping.GenerateSource(mcVersion, mappingType)
			return err
		})
		c.String(http.StatusAccepted, "Started decompiling, please wait")
	})
	g.GET("/api/source/get", func(c *gin.Context) {
		mcVersion, mappingType, class := c.Query("version"), c.Query("type"), c.Query("class")
		if mcVersion == "" || mappingType == "" || class == "" {
			c.String(http.StatusBadRequest, "Missing query parameter(s)")
			return
		}
		if !mapping.IsAvailable(mcVersion, mappingType) {
			c.String(http.StatusPreconditionFailed, "Use /load before getting")
			return
		}
		path := global.GetSourceFolder(global.NamedImpl{Name: mappingType}, mcVersion)
		targetPath := filepath.Join(path, class+".java")
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			c.String(http.StatusNotFound, "%s not found ", targetPath)
		}
		content, err := os.ReadFile(targetPath)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read file")
			return
		}
		c.Data(http.StatusOK, "text/plain; charset=utf-8", content)
	})
}
