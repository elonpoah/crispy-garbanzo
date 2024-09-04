package utils

import (
	"crispy-garbanzo/global"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

type JwtClaims struct {
	Aid int    `json:"aid"`
	Src string `json:"src"`
	Iat int64  `json:"iat"`
	Exp int64  `json:"exp"`
}

func CreateToken(aid int, src string, key string, exp int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aid": aid,
		"src": src,
		"iat": time.Now().Unix(),
		"exp": time.Now().Unix() + int64(exp),
	})

	tokenString, err := token.SignedString([]byte(key))
	if nil != err {
		global.FPG_LOG.Error("signing error: %s", zap.Error(err))
		return ""
	}
	return tokenString
}

func ParseToken(tokenString string, key string) (JwtClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(key), nil
	})

	var jwtClaims JwtClaims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if err := mapstructure.Decode(claims, &jwtClaims); err != nil {
			return jwtClaims, fmt.Errorf("signing invalid: %v", err)
		} else {
			return jwtClaims, nil
		}
	} else {
		return jwtClaims, fmt.Errorf("signing invalid: %v", err)
	}
}

func GetUserID(c *gin.Context) (uid int, err error) {
	username, exists := c.Get("uid")
	if !exists {
		return 0, errors.New("user ID is invalid")
	}
	usernameStr, ok := username.(int)
	if !ok {
		return 0, errors.New("user ID is invalid")
	}
	return usernameStr, nil
}
