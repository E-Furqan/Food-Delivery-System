package database

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

type RepositoryInterface interface {
	CreateRole(role *model.Role) error
	GetRole(RoleId uint, role *model.Role) error
	GetAllRoles() ([]model.Role, error)
	BulkCreateRoles(roles []model.Role) error
	DeleteRole(role *model.Role) error

	CreateUser(user *model.User) error
	LoadUserWithRoles(user *model.User) error
	GetUser(columnName string, findParameter interface{}, user *model.User) error
	FetchUsersWithRoles(columnName string, order string) ([]model.User, error)
	Update(user *model.User, update_user *model.User) error
	DeleteUser(user *model.User) error
	DeleteUserRoleInfo(ID uint, columnName string) error

	AddUserRole(userId uint, roleId uint) error
	UpdateUserActiveRole(user *model.User) error
	UpdateRoleStatus(user *model.User) error
}
