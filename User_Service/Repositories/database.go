package database

import (
	"fmt"
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/models"
	"gorm.io/gorm"
)

// Repository struct to handle dependency injection
type Repository struct {
	DB *gorm.DB
}

// NewRepository is a constructor function to initialize the repository with a DB connection
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

// CreateRole inserts a new role into the database
func (repo *Repository) CreateRole(role *model.Role) error {
	result := repo.DB.Create(role)
	if result.Error != nil {
		log.Printf("Error creating role: %v", result.Error)
		return result.Error
	}
	return nil
}

// CreateUser inserts a new user into the database
func (repo *Repository) CreateUser(user *model.User) error {
	result := repo.DB.Create(user)
	repo.LoadUserWithRoles(user.UserId)
	return result.Error
}

// LoadUserWithRole loads a user with its associated role from the database
func (repo *Repository) LoadUserWithRoles(userID uint) (model.User, error) {
	var userWithRoles model.User
	// Load the user with the associated roles
	err := repo.DB.Preload("Roles").First(&userWithRoles, userID).Error
	if err != nil {
		log.Printf("Error loading user with roles: %v", err)
		return model.User{}, err
	}

	log.Printf("Successfully loaded user with roles: %v", userWithRoles.Roles)
	return userWithRoles, nil
}

// FindUser retrieves a user
func (repo *Repository) FindUser(columnName string, findParameter interface{}, user *model.User) error {
	query := fmt.Sprintf("%s = ?", columnName)
	err := repo.DB.Where(query, findParameter).First(user).Error

	// Load roles for the user
	userWithRoles, err := repo.LoadUserWithRoles(user.UserId)
	if err != nil {
		return err // Return error if loading roles fails
	}

	// Update the user with roles
	*user = userWithRoles

	return err
}

// WhereRoleID retrieves a role by user role ID
func (repo *Repository) FindRole(RoleId interface{}, role *model.Role) error {

	err := repo.DB.Where("role_id = ?", RoleId).First(role).Error
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

// func (repo *Repository) Save(user *model.User) error {
// 	return repo.DB.Save(user).Error
// }

func (repo *Repository) Update(user *model.User, update_user *model.User) error {
	err := repo.DB.Model(user).Updates(update_user).Error
	log.Print("inside update")
	log.Print(update_user.Roles)
	log.Print(user.Roles)
	repo.LoadUserWithRoles(user.UserId)
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
	// Delete all entries in the user_roles table that reference the given roleId
	result := repo.DB.Where(query, ID).Delete(&UserRole)

	return result.Error
}

//	func (repo *Repository) DeleteUserRoleInfo(roleId uint, tablename string) error {
//		var UserRole model.UserRole
//		query := fmt.Sprintf("%s = ?", tablename)
//		// Delete all entries in the user_roles table that reference the given roleId
//		result := repo.DB.Where(query, roleId).Delete(&UserRole)
//		return result.Error
//	}
func (repo *Repository) AddUserRole(userId uint, roleId uint) error {

	var existingUserRole model.UserRole
	if err := repo.DB.Where("user_user_id = ? AND role_role_id = ?", userId, roleId).First(&existingUserRole).Error; err == nil {
		// Role already exists for this user, so return nil or a custom message
		return nil // or return fmt.Errorf("role already exists for this user")
	}

	userRole := model.UserRole{
		UserId: userId,
		RoleId: roleId,
	}
	err := repo.DB.Create(&userRole).Error
	return err
}