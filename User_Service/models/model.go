package model

type User struct {
	UserId   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName string `gorm:"size:100;not null"`
	Username string `gorm:"size:100;not null;uniqueIndex"`
	Password string `gorm:"size:100;not null"`

	Email       string `gorm:"size:100;uniqueIndex"`
	PhoneNumber uint   `gorm:"size:100;uniqueIndex"`
	Address     string `gorm:"size:100"`

	RoleId string `gorm:"size:100;not null"`
	Role   Role   `gorm:"foreignKey:Role_id;references:Role_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	RoleStatus string `gorm:"size:100"`
}

type Role struct {
	RoleId   string `gorm:"primaryKey;size:100"`
	RoleType string `gorm:"size:100"`
}
