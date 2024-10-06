package middleware

import (
	"crispy-garbanzo/common/response"
	"crispy-garbanzo/global"
	"crispy-garbanzo/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if len(token) == 0 {
			response.NoAuth(response.TokenMissing, c)
			c.Abort()
			return
		}

		jwtClaims, err := utils.ParseToken(token, global.FPG_CONFIG.Jwt.Key)

		if err != nil {
			global.FPG_LOG.Error("签名错误或者过期!", zap.Error(err))
			response.NoAuth(response.TokenError, c)
			c.Abort()
			return
		}

		if jwtClaims.Exp < time.Now().Unix() {
			response.NoAuth(response.TokenExpired, c)
			c.Abort()
			return
		}
		c.Set("uid", jwtClaims.Aid)

		c.Next()
	}
}
