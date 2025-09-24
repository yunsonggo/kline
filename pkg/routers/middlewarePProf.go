package routers

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/yunsonggo/kline/pkg/xconfig"
)

type PProfMiddleware struct {
	// []string{"127.0.0.1", "192.168.1.0/24", "192.168.2.0/24"}
	AllowOrigins []string `json:"allowOrigins"`
}

func NewPProfMiddleware(c xconfig.GinConf) *PProfMiddleware {
	pm := &PProfMiddleware{}
	if len(c.AllowPProfOrigins) > 0 {
		pm.AllowOrigins = c.AllowPProfOrigins
	}
	return pm
}

func (pm *PProfMiddleware) Register(r *gin.Engine) {
	authorized := r.Group("/debug/pprof", func(c *gin.Context) {
		clientIP := c.ClientIP()
		if len(pm.AllowOrigins) > 0 {
			for _, ip := range pm.AllowOrigins {
				if ip == clientIP {
					c.Next()
					return
				}
			}
		}
		c.AbortWithStatus(403)
	})
	pprof.Register(authorized)
}
