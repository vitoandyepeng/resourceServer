package middleware

import (
	"common/utils"
	"crypto/md5"
	"data"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {

			c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Cookie, Accept, Authorization")

			c.AbortWithStatus(http.StatusOK)
		} else {
			c.Next()
		}
	}
}

func LoginRequire() gin.HandlerFunc {

	return func(c *gin.Context) {
		var req data.Req
		token, ok := c.GetPostForm("md5")
		var key string
		if ok {
			key = data.Config.PrivateKey
		} else {
			err := c.Bind(&req)
			if err != nil {
				utils.WErr("LoginRequire bind err.", err.Error())
				Echo(c, http.StatusUnauthorized, "")
				c.Abort()
				return
			}
			token = req.Md5
			key = fmt.Sprintf("%s%d", data.Config.PrivateKey, req.Id)
			c.Set("data", req)
		}

		m := md5.New()
		m.Write([]byte(key))
		cipherStr := m.Sum(nil)
		str := hex.EncodeToString(cipherStr)
		if str != token {
			utils.WErr("LoginRequire md5 err.")
			Echo(c, http.StatusUnauthorized, "")
			c.Abort()
			return
		}

		c.Next()
	}
}

func Echo(c *gin.Context, code int, str string)  {
	c.Writer.Header().Set("Content-type", "application/text")
	c.Writer.WriteHeader(code)
	n, err := c.Writer.Write([]byte(str))
	if err != nil {
		utils.WErr("Echo err", n, err.Error())
	}
}