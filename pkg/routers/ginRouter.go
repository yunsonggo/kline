package routers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/yunsonggo/kline/pkg/xconfig"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewGinRouter(c xconfig.GinConf) *gin.Engine {
	switch c.Mode {
	case gin.DebugMode:
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
		gin.DefaultWriter = io.Discard
	}
	router := gin.New()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	if c.UseCors {
		corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, c.AllowHeaders...)
		corsConfig.AllowOrigins = append(corsConfig.AllowOrigins, c.AllowOrigins...)
	} else {
		corsConfig.AllowHeaders = []string{"*"}
		corsConfig.AllowAllOrigins = true
	}
	router.Use(cors.New(corsConfig))

	if c.UsePProf {
		pm := NewPProfMiddleware(c)
		pm.Register(router)
	}
	router.StaticFS(c.StaticPrefix, http.Dir(c.StaticPath))

	router.NoRoute(func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"code":  http.StatusNotFound,
			"msg":   "Bad request",
			"error": fmt.Sprintf("Method %s URI %s,not found", ctx.Request.Method, ctx.Request.RequestURI),
		})
	})
	
	return router
}
