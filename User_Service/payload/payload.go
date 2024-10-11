package payload

type Input struct {
	Username string `json:"username_admin"`
	Password string `json:"password_admin"`
	Role_id  string `json:"role_id"`
}
