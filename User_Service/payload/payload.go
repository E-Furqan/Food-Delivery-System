package payload

type Input struct {
	Username string `json:"usernameAdmin"`
	Password string `json:"passwordAdmin"`
	RoleId   string `json:"roleId"`
}

type Order struct {
	ColumnName string `json:"columnName"`
	OrderType  string `json:"OrderType"`
}
