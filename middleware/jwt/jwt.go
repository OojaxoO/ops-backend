package jwt

import (
	"net/http"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"ops-backend/pkg/e"
	"ops-backend/pkg/util"
)

type Header struct {
	Token 		string `header:"Authorization"`
}

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token, err := GetToken(c)
		if err != nil {
			code = e.INVALID_PARAMS
		}
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			_, err := util.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				default:
					code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
				}
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}

func GetToken(c *gin.Context) (string, error) {
	h := Header{}
	if err := c.ShouldBindHeader(&h); err != nil {
		return "", err
	}
	token := h.Token
	tokenSplit := strings.Split(token, "Bearer ")
	if len(tokenSplit) < 2 {
		err := fmt.Errorf("token格式错误")
		return "", err 
	}
	token = tokenSplit[1]
	return token, nil
}
