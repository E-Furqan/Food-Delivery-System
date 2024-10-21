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
	ColumnName string `json:"column_name"`
	OrderType  string `json:"order-type"`
}

type RoleSwitch struct {
	NewRoleID uint `json:"switch_to"`
}

var RolesList = []struct {
	RoleId   uint
	RoleType string
}{
	{1, "Customer"},
	{2, "Delivery driver"},
	{3, "Admin"},
}

type ProcessOrder struct {
	OrderId         uint   `json:"order_id"`
	UserID          uint   `json:"user_id"`
	DeliverDriverID uint   `json:"delivery_driver"`
	OrderStatus     string `json:"order_status"`
}

func GetOrderTransitions() map[string]string {
	return map[string]string{
		"Waiting For Delivery Driver": "In for delivery",
		"In for delivery":             "Delivered",
		"Delivered":                   "Completed",
	}
}
