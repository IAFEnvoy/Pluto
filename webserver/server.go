package webserver

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"pluto/global"
	"pluto/util"
	"strconv"
)

func Launch() error {
	gin.DefaultWriter = &util.SlogWriter{Level: slog.LevelInfo}
	gin.DefaultErrorWriter = &util.SlogWriter{Level: slog.LevelError}
	g := gin.Default()
	g.Use(cors.New(global.Config.Cors))
	g.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "{\"version\":\"%s\"}", global.Version)
	})
	initMappingApis(g)
	initSourceApi(g)
	err := g.Run(":" + strconv.Itoa(global.Config.Port))
	return err
}
