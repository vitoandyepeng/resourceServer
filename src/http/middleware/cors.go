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

		err := c.Bind(&req)
		if err != nil {
			utils.WErr("LoginRequire bind err.", err.Error())
			Echo(c, http.StatusUnauthorized, "")
			c.Abort()
			return
		}

		//if "abcdefg" != req.Md5 {
		//	utils.WErr("LoginRequire md5 err.", err.Error())
		//	global.Echo(c, http.StatusUnauthorized, res)
		//	c.Abort()
		//	return
		//}

		m := md5.New()
		m.Write([]byte(fmt.Sprintf("%s%d", data.Config.PrivateKey, req.Id)))
		cipherStr := m.Sum(nil)
		str := hex.EncodeToString(cipherStr)
		if str != req.Md5 {
			utils.WErr("LoginRequire md5 err.", err.Error())
			Echo(c, http.StatusUnauthorized, "")
			c.Abort()
			return
		}

		c.Set("data", req)

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