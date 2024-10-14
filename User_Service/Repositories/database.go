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
	repo.LoadUserWithRole(user.UserId)
	return result.Error
}

// LoadUserWithRole loads a user with its associated role from the database
func (repo *Repository) LoadUserWithRole(userID uint) (model.User, error) {
	var userWithRole model.User
	// Load the user with the associated role
	err := repo.DB.Preload("Role").First(&userWithRole, userID).Error
	if err != nil {
		log.Printf("Error loading user with role: %v  %s", err, userWithRole.Role.RoleType)
		return model.User{}, err // Return an empty User struct and the error
	}
	log.Printf("Error loading user with role: %v  %s", err, userWithRole.Role.RoleType)
	return userWithRole, err // Return the loaded user and a nil error
}

// WhereUsername retrieves a user by username
func (repo *Repository) FindUser(FindParameter string, user *model.User) error {
	err := repo.DB.Where("username = ?", FindParameter).First(user).Error
	return err
}

// WhereRoleID retrieves a role by user role ID
func (repo *Repository) FindRole(FindParameter string, role *model.Role) error {

	err := repo.DB.Where("RoleId = ?", FindParameter).First(role).Error
	return err
}

func (repo *Repository) PreloadInOrder(columnName string, order string) ([]model.User, error) {

	if columnName == "" {
		columnName = "UserId"
	}
	if order == "" {
		order = "asc"
	}
	query := fmt.Sprintf("%s %v", columnName, order)

	var user_data []model.User
	err := repo.DB.Preload("Role").Order(query).Find(&user_data).Error

	return user_data, err
}

func (repo *Repository) RoleInOrder(columnName string, order string) ([]model.Role, error) {

	if columnName == "" {
		columnName = "RoleId"
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
	repo.LoadUserWithRole(user.UserId)
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
