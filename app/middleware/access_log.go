package middleware

import (
	"bytes"
	"gin-web/app/vars"
	"gin-web/libraty/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type AccessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

var logExcludeApis = []string{
	"api/upload", "api/login", "api/logout", "api/profile",
}

func (w AccessLogWriter) Write(p []byte) (int, error) {
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 会影响文件上传参数接收，文件上传不写日志
		if vars.YmlConfig.GetBool("Debug") || inExclude(c.Request.RequestURI) {
			c.Next()
		} else {
			beginTime := time.Now()
			c.Next()
			endTime := time.Now()
			user, _ := c.Get(vars.LoginUserKey)

			// 请求
			//var data interface{}
			//_ = c.ShouldBindBodyWith(&data, binding.JSON)
			// 响应
			//bodyWriter := &AccessLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			//c.Writer = bodyWriter
			//bodyWriter.body.String()

			logger.NewLog(logrus.InfoLevel, "access-log").Log(logrus.Fields{
				"uid":   user.(*vars.LoginUser).UserId,
				"uname": user.(*vars.LoginUser).Username,
				"uri":   c.Request.RequestURI,
				"ip":    c.ClientIP(),
				"begin": beginTime.Format(vars.DateTimeFormat),
				"end":   endTime.Format(vars.DateTimeFormat),
			}, "")
		}
	}
}

func inExclude(uri string) bool {
	for _, i2 := range logExcludeApis {
		if strings.Contains(uri, i2) {
			return true
		}
	}
	return false
}
