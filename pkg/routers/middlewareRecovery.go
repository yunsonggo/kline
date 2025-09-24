package routers

import (
	"net"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yunsonggo/kline/pkg/xloger"
)

type RecoveryMiddleware struct {
	Logger xloger.Logger
}

func NewMRecoveryMiddleware(logger xloger.Logger) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		Logger: logger,
	}
}

func (m *RecoveryMiddleware) NewHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				isBroken := m.isBrokenErr(panicErr)
				req, _ := httputil.DumpRequest(ctx.Request, true)
				m.Logger.Error(ctx, "PANIC",
					xloger.Field{
						Key:   "Request",
						Value: req,
					}, xloger.Field{
						Key:   "Error",
						Value: panicErr,
					}, xloger.Field{
						Key:   "Stack",
						Value: debug.Stack(),
					})
				if isBroken {
					ctx.Error(panicErr.(error))
					ctx.Abort()
				}
			}
		}()
		ctx.Next()
	}
}

func (m *RecoveryMiddleware) isBrokenErr(panicErr any) bool {
	netErr, ok := panicErr.(*net.OpError)
	if !ok {
		return false
	}
	sysErr, ok := netErr.Err.(*os.SyscallError)
	if !ok {
		return false
	}
	message := strings.ToLower(sysErr.Error())
	if strings.Contains(message, "broken pipe") || strings.Contains(message, "connection reset by peer") {
		return true
	}
	return false
}
