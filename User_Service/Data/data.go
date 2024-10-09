package data

type User struct {
	User_id   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Full_Name string `gorm:"size:100;not null"`
	Username  string `gorm:"size:100;not null"`
	Password  string `gorm:"size:100;not null"`

	Email        string `gorm:"size:100"`
	Phone_number uint   `gorm:"size:100"`
	Address      string `gorm:"size:100"`

	Role_id string `gorm:"size:100;not null"`
	Role    Role   `gorm:"foreignKey:Role_id;references:Role_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Role_status string `gorm:"size:100"`
}

type Role struct {
	Role_id   string `gorm:"primaryKey;size:100"`
	Role_type string `gorm:"size:100"`
}
