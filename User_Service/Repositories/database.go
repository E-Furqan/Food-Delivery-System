package database

import (
	"fmt"
	"log"

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

func (repo *Repository) CreateRole(role *model.Role) error {
	result := repo.DB.Create(role)
	if result.Error != nil {
		log.Printf("Error creating role: %v", result.Error)
		return result.Error
	}
	return nil
}

func (repo *Repository) CreateUser(user *model.User) error {
	result := repo.DB.Create(user)
	log.Print(user)
	repo.LoadUserWithRoles(user)
	log.Print(user)
	return result.Error
}

func (repo *Repository) LoadUserWithRoles(user *model.User) error {

	log.Printf("Successfully loaded user with roles: %v", user.Roles)

	err := repo.DB.Preload("Roles").First(user, user.UserId).Error
	if err != nil {
		log.Printf("Error loading user with roles: %v", err)
		return err
	}

	log.Printf("Successfully loaded user with roles: %v", user.Roles)
	return nil
}

func (repo *Repository) GetUser(columnName string, findParameter interface{}, user *model.User) error {
	query := fmt.Sprintf("%s = ?", columnName)
	err := repo.DB.Where(query, findParameter).First(user).Error
	if err != nil {
		log.Printf("Error : %s", err)
		return err
	}
	err = repo.LoadUserWithRoles(user)

	return err
}

func (repo *Repository) GetRole(RoleId uint, role *model.Role) error {

	err := repo.DB.Where("role_id = ?", RoleId).First(role).Error
	return err
}

func (repo *Repository) GetAdminRole(RoleIds []uint) error {
	var role *model.Role
	err := repo.DB.Where("role_id IN ? AND role_type = ?", RoleIds, "admin").First(role).Error
	return err
}

func (repo *Repository) PreloadInOrder(columnName string, order string) ([]model.User, error) {

	if columnName == "" {
		columnName = "user_id"
	}
	if order == "" {
		order = "asc"
	}
	query := fmt.Sprintf("%s %v", columnName, order)

	var user_data []model.User
	err := repo.DB.Preload("Roles").Order(query).Find(&user_data).Error

	return user_data, err
}

func (repo *Repository) RoleInOrder(columnName string, order string) ([]model.Role, error) {

	if columnName == "" {
		columnName = "role_id"
	}
	if order == "" {
		order = "asc"
	}

	var user_data []model.Role

	query := fmt.Sprintf("%s %v", columnName, order)

	err := repo.DB.Order(query).Find(&user_data).Error

	return user_data, err
}

func (repo *Repository) Update(user *model.User, update_user *model.User) error {
	err := repo.DB.Model(user).Updates(update_user).Error
	_ = repo.LoadUserWithRoles(user)
	return err
}

func (repo *Repository) DeleteUser(user *model.User) error {
	err := repo.DB.Delete(user).Error
	return err
}

func (repo *Repository) DeleteRole(Role *model.Role) error {
	err := repo.DB.Delete(Role).Error
	return err
}

func (repo *Repository) CheckUserExistence(username, email string, phoneNumber int) (bool, error) {
	var count int64

	err := repo.DB.Model(&model.User{}).Where("username = ? OR phone_number = ? OR email = ?", username, phoneNumber, email).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *Repository) BulkCreateRoles(roles []model.Role) error {
	result := repo.DB.Create(&roles)
	if result.Error != nil {
		log.Printf("Error creating role: %v", result.Error)
		return result.Error
	}
	return nil
}

func (repo *Repository) DeleteUserRoleInfo(ID uint, columnName string) error {
	var UserRole model.UserRole
	query := fmt.Sprintf("%s = ?", columnName)

	result := repo.DB.Where(query, ID).Delete(&UserRole)

	return result.Error
}

func (repo *Repository) AddUserRole(userId uint, roleId uint) error {

	var existingUserRole model.UserRole
	if err := repo.DB.Where("user_user_id = ? AND role_role_id = ?", userId, roleId).First(&existingUserRole).Error; err == nil {
		return nil
	}

	userRole := model.UserRole{
		UserId: userId,
		RoleId: roleId,
	}
	err := repo.DB.Create(&userRole).Error
	return err
}

func (repo *Repository) UpdateUserActiveRole(user *model.User) error {
	return repo.DB.Model(user).Update("ActiveRole", user.ActiveRole).Error
}

func (repo *Repository) GetDeliveryDrivers() (model.User, error) {

	tx := repo.DB.Begin()
	var driver model.User

	err := repo.DB.Where("active_role = ? AND (role_status = ? OR role_status = ?)",
		"Delivery driver", "not active", "not available").First(&driver).Error

	if err != nil {
		tx.Rollback()
		return driver, err
	}

	err = repo.DB.Model(&driver).Update("roleStatus", "not available").Error

	if err != nil {
		tx.Rollback()
		return driver, err
	}

	return driver, nil
}
