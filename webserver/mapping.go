package webserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pluto/global"
	"pluto/mapping"
)

func initMappingApis(g *gin.Engine) {
	g.GET("/api/mapping/load", func(c *gin.Context) {
		mcVersion, mappingType := c.Query("version"), c.Query("type")
		if mcVersion == "" || mappingType == "" {
			c.String(http.StatusBadRequest, "Missing query parameter(s)")
			return
		}
		if mapping.CachedMapping(mcVersion, mappingType) {
			c.String(http.StatusOK, "Loaded", global.Version)
			return
		}
		_, err := mapping.LoadMapping(mcVersion, mappingType)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.String(http.StatusCreated, "Loading Complete")
	})
	g.GET("/api/mapping/search", func(c *gin.Context) {
		mcVersion, mappingType, keyword := c.Query("version"), c.Query("type"), c.Query("keyword")
		if mcVersion == "" || mappingType == "" || keyword == "" {
			c.String(http.StatusBadRequest, "Missing query parameter(s)")
			return
		}
		if len(keyword) <= 2 {
			c.String(http.StatusBadRequest, "Keyword must contain at least three characters")
			return
		}
		if !mapping.CachedMapping(mcVersion, mappingType) {
			c.String(http.StatusPreconditionFailed, "Use /load before searching")
			return
		}
		mappings, err := mapping.LoadMapping(mcVersion, mappingType)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, mappings.Search(keyword, 20))
	})
}
