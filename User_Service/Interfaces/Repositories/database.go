package database

import (
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
	return result.Error
}

// LoadUserWithRole loads a user with its associated role from the database
func (repo *Repository) LoadUserWithRole(userID uint) (model.User, error) {
	var userWithRole model.User

	// Load the user with the associated role
	err := repo.DB.Preload("Role").First(&userWithRole, userID).Error
	if err != nil {
		log.Printf("Error loading user with role: %v  %s", err, userWithRole.Role.Role_type)
		return model.User{}, err // Return an empty User struct and the error
	}
	log.Printf("Error loading user with role: %v  %s", err, userWithRole.Role.Role_type)
	return userWithRole, err // Return the loaded user and a nil error
}

// WhereUsername retrieves a user by username
func (repo *Repository) Find_User_By_Username(username string, user *model.User) error {
	err := repo.DB.Where("username = ?", username).First(user).Error
	return err
}

// WhereRoleID retrieves a role by user role ID
func (repo *Repository) Find_Role_By_Role_Id(Role_id string, role *model.Role) error {

	err := repo.DB.Where("role_id = ?", Role_id).First(role).Error
	return err
}

func (repo *Repository) Preload_in_order() ([]model.User, error) {
	var user_data []model.User
	err := repo.DB.Preload("Role").Order("User_id asc").Find(&user_data).Error
	return user_data, err
}

func (repo *Repository) Role_in_Asc_order() ([]model.Role, error) {
	var user_data []model.Role
	err := repo.DB.Order("Role_id asc").Find(&user_data).Error
	return user_data, err
}

func (repo *Repository) Save(user *model.User) error {
	return repo.DB.Save(user).Error
}

func (repo *Repository) Update(user *model.User, update_user *model.User) error {
	err := repo.DB.Model(user).Updates(update_user).Error
	return err
}

func (repo *Repository) Preload_Role_first(userWithRole *model.User, User_id int) error {
	err := repo.DB.Preload("Role").First(userWithRole, User_id).Error
	return err
}

func (repo *Repository) Delete_user(user *model.User) error {
	err := repo.DB.Delete(user).Error
	return err
}

func (repo *Repository) Delete_role(Role *model.Role) error {
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
