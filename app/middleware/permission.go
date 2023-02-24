package middleware

import (
	"gin-web/app/response"
	"gin-web/app/vars"
	"github.com/gin-gonic/gin"
	"strconv"
)

func CheckPermission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if user, exists := ctx.Get(vars.LoginUserKey); !exists {
			response.Fail(ctx, "登录信息获取失败")
			return
		} else {
			// 超级管理员角色不需要检查权限
			if user.(*vars.LoginUser).RoleId == 1 {
				ctx.Next()
			} else {
				requestUrl := urlParse(ctx.Request.URL.Path) // 路由例如 /admin/user/index
				method := ctx.Request.Method        // 方法 GET/POST/PUT...
				// 判断所拥有的的角色是否具有某个权限即可
				ok, err := vars.Casbin.Enforce(strconv.Itoa(int(user.(*vars.LoginUser).RoleId)), requestUrl, method)
				if err != nil {
					response.Fail(ctx, "权限验证出错了:"+err.Error())
					return
				}

				if !ok {
					response.NoPermission(ctx)
					return
				} else {
					ctx.Next()
				}
			}
		}
	}
}

func urlParse(u string) string {
	return u[len(vars.ApiPrefix)+1:]
}
