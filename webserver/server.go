package webserver

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"pluto/global"
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
	r.GET("/api/")
	err := r.Run(":" + strconv.Itoa(global.Config.Port))
	return err
}
