package payload

type Role struct {
	RoleType string `json:"roleType"`
	RoleId   uint   `json:"roleId"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserSearch struct {
	ColumnName      string      `json:"columnName"`
	SearchParameter interface{} `json:"searchParameter"`
}

type Order struct {
	ColumnName string `json:"columnName"`
	OrderType  string `json:"OrderType"`
}

type RoleSwitch struct {
	NewRoleID uint `json:"SwitchTo"`
}

var RolesList = []struct {
	RoleId   uint
	RoleType string
}{
	{1, "Customer"},
	{2, "Delivery driver"},
	{3, "Admin"},
}
