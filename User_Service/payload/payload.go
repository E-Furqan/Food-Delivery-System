package payload

type Input struct {
	Username string `json:"username_admin" binding:"required"`
	Password string `json:"password_admin" binding:"required"`
	Role_id  string `json:"role_id" binding:"required"`
}
