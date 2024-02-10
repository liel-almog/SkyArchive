package configs

type CustomJwtClaims struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Id       int64  `json:"id"`
}
