package vars

type LoginUser struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	RoleId   int64  `json:"role_id"`
}
