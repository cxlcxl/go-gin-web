package errs

import "errors"

var (
	LoginFailToken = errors.New("登陆失败，TOKEN 创建失败")
	LoginUserVoid  = errors.New("账号已失效不可登陆")
	LoginPassErr   = errors.New("密码错误")
)
