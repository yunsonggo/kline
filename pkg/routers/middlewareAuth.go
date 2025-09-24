package routers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yunsonggo/kline/pkg/xjwt"
)

type CheckRole func(id string) bool
type AuthMiddleware struct {
	Jwt       xjwt.Jwt
	Header    string
	UseBearer bool
	UseRole   bool
	CheckRole CheckRole
}

func NewAuthMiddleware(jwt xjwt.Jwt, header string, useBearer, useRole bool, roleFunc CheckRole) *AuthMiddleware {
	return &AuthMiddleware{
		Jwt:       jwt,
		Header:    header,
		UseBearer: useBearer,
		UseRole:   useRole,
		CheckRole: roleFunc,
	}
}

func (am *AuthMiddleware) NewHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get(am.Header)
		if authHeader == "" {
			ctx.JSON(http.StatusOK, gin.H{
				"code":  http.StatusUnauthorized,
				"msg":   "Bad request",
				"error": "unauthenticated",
			})
			ctx.Abort()
			return
		}
		token := ""
		if am.UseBearer {
			parts := strings.SplitN(authHeader, " ", 2)
			if !(len(parts) == 2 && parts[0] == "Bearer") {
				ctx.JSON(http.StatusOK, gin.H{
					"code":  http.StatusUnauthorized,
					"msg":   "Bad request",
					"error": "unauthenticated",
				})
				ctx.Abort()
				return
			}
			token = parts[1]
		} else {
			token = authHeader
		}
		claims, err := am.Jwt.Parse(token)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code":  http.StatusUnauthorized,
				"msg":   "Bad request",
				"error": "unauthenticated",
			})
			ctx.Abort()
			return
		}
		if am.UseRole {
			idStr, ok := claims["id"]
			if !ok {
				ctx.JSON(http.StatusOK, gin.H{
					"code":  http.StatusUnauthorized,
					"msg":   "Bad request",
					"error": "unauthenticated",
				})
				ctx.Abort()
				return
			}
			id, ok := idStr.(string)
			if !ok {
				ctx.JSON(http.StatusOK, gin.H{
					"code":  http.StatusUnauthorized,
					"msg":   "Bad request",
					"error": "unauthenticated",
				})
				ctx.Abort()
				return
			}
			if !am.CheckRole(id) {
				ctx.JSON(http.StatusOK, gin.H{
					"code":  http.StatusUnauthorized,
					"msg":   "Bad request",
					"error": "unauthenticated",
				})
				ctx.Abort()
				return
			}
		}
		ctx.Set("claims", claims)
		ctx.Next()
	}
}
