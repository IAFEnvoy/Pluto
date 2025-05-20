package webserver

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"pluto/global"
	"pluto/mapping"
	"pluto/util"
	"strconv"
)

func Launch() error {
	gin.DefaultWriter = &util.SlogWriter{Level: slog.LevelInfo}
	gin.DefaultErrorWriter = &util.SlogWriter{Level: slog.LevelError}
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "{\"version\":\"%s\"}", global.Version)
	})
	r.GET("/api/mapping/load", func(c *gin.Context) {
		mcVersion, mappingType := c.Query("version"), c.Query("type")
		if mcVersion == "" || mappingType == "" {
			c.String(400, "Missing query parameter(s)")
			return
		}
		if mapping.CachedMapping(mcVersion, mappingType) {
			c.String(200, "Loaded", global.Version)
			return
		}
		_, err := mapping.LoadMapping(mcVersion, mappingType)
		if err != nil {
			c.String(400, err.Error())
			return
		}
		c.String(201, "Loading Complete")
	})
	r.GET("/api/mapping/search", func(c *gin.Context) {
		mcVersion, mappingType, keyword := c.Query("version"), c.Query("type"), c.Query("keyword")
		if mcVersion == "" || mappingType == "" || keyword == "" {
			c.String(400, "Missing query parameter(s)")
			return
		}
		if len(keyword) <= 2 {
			c.String(400, "Keyword must contain at least three characters")
			return
		}
		if !mapping.CachedMapping(mcVersion, mappingType) {
			c.String(412, "Use /load before searching")
			return
		}
		mappings, err := mapping.LoadMapping(mcVersion, mappingType)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		c.JSON(200, mappings.Search(keyword, 20))
	})
	err := r.Run(":" + strconv.Itoa(global.Config.Port))
	return err
}
