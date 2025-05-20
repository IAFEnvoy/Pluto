package webserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pluto/mapping"
	"time"
)

func initMappingApis(g *gin.Engine) {
	g.GET("/api/mapping/load", RateLimiterMiddleware(5*time.Second, 2), func(c *gin.Context) {
		mcVersion, mappingType := c.Query("version"), c.Query("type")
		if mcVersion == "" || mappingType == "" {
			c.String(http.StatusBadRequest, "Missing query parameter(s)")
			return
		}
		if mapping.CachedMapping(mcVersion, mappingType) {
			c.String(http.StatusOK, "Loaded")
			return
		}
		_, err := mapping.LoadMapping(mcVersion, mappingType)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.String(http.StatusCreated, "Loading Complete")
	})
	g.GET("/api/mapping/search", RateLimiterMiddleware(2*time.Second, 5), func(c *gin.Context) {
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
