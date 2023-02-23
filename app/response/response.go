package response

import (
	"fmt"
	"gin-web/app/vars"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ReturnJson ...
func ReturnJson(c *gin.Context, httpCode int, dataCode int, msg string, data interface{}) {
	c.JSON(httpCode, gin.H{
		"code":    dataCode,
		"message": msg,
		"data":    data,
	})
	c.Abort()
}

// Html ...
func Html(c *gin.Context, page string, data interface{}) {
	c.HTML(http.StatusOK, fmt.Sprintf("%s.tmpl", page), data)
	//c.Abort()
}

// Success 直接返回成功
func Success(c *gin.Context, data interface{}) {
	ReturnJson(c, http.StatusOK, vars.ResponseCodeOk, vars.ResponseMsg[vars.ResponseCodeOk], data)
}

// NoPermission 没有权限
func NoPermission(c *gin.Context) {
	ReturnJson(c, 419, 419, "权限不足，不可访问", "")
}

// Fail 失败的业务逻辑
func Fail(c *gin.Context, msg string) {
	ReturnJson(c, http.StatusBadRequest, 400, msg, nil)
}

// TokenExpired Token 过期
func TokenExpired(c *gin.Context) {
	ReturnJson(c, http.StatusBadRequest, 418, "token 过期，请手动刷新页面重新登陆", "")
}
