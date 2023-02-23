package types

type Login struct {
	Email string `json:"email" binding:"required,email"`
	Pass  string `json:"pass" binding:"required,pass"`
}
