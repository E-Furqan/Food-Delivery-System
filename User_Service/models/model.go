package model

type User struct {
	UserId      uint   `gorm:"primaryKey;column:user_id;autoIncrement" json:"user_id"`
	FullName    string `gorm:"size:100;not null;column:full_name" json:"fullName"`
	Username    string `gorm:"size:100;not null;uniqueIndex;column:username" json:"username"`
	Password    string `gorm:"size:100;not null;column:password" json:"password"`
	Email       string `gorm:"size:100;uniqueIndex;column:email" json:"email"`
	PhoneNumber uint   `gorm:"uniqueIndex;column:phone_number" json:"phoneNumber"`
	Address     string `gorm:"size:100;column:address" json:"address"`
	RoleStatus  string `gorm:"size:100;column:role_status" json:"roleStatus"`
	ActiveRole  string `gorm:"size:100;column:active_role" json:"activeRole"`
	Roles       []Role `gorm:"many2many:user_roles;foreignKey:user_id;joinForeignKey:user_user_id;References:role_id;joinReferences:role_role_id" json:"roles"` // Establish many-to-many relationship
}

type Role struct {
	RoleId   uint   `gorm:"primaryKey;column:role_id;autoIncrement" json:"roleId"`
	RoleType string `gorm:"size:100;column:role_type" json:"roleType"`
}

type UserRole struct {
	UserId uint `gorm:"primaryKey;column:user_user_id" json:"userId"`
	RoleId uint `gorm:"primaryKey;column:role_role_id" json:"roleId"`
}
