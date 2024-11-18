package model

type Role struct {
	RoleId   uint   `gorm:"primaryKey;column:role_id" json:"roleId"`
	RoleType string `gorm:"size:100;column:role_type" json:"roleType"`
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
