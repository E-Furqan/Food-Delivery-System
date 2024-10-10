package database

import (
	"log"

	entity "github.com/E-Furqan/Food-Delivery-System/Entity"
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
func (repo *Repository) CreateRole(role *entity.Role) error {
	result := repo.DB.Create(role)
	if result.Error != nil {
		log.Printf("Error creating role: %v", result.Error)
		return result.Error
	}
	return nil
}

// CreateUser inserts a new user into the database
func (repo *Repository) CreateUser(user *entity.User) error {
	result := repo.DB.Create(user)
	return result.Error
}

// LoadUserWithRole loads a user with its associated role from the database
func (repo *Repository) LoadUserWithRole(userID uint) (entity.User, error) {
	var userWithRole entity.User

	// Load the user with the associated role
	err := repo.DB.Preload("Role").First(&userWithRole, userID).Error
	if err != nil {
		log.Printf("Error loading user with role: %v  %s", err, userWithRole.Role.Role_type)
		return entity.User{}, err // Return an empty User struct and the error
	}
	log.Printf("Error loading user with role: %v  %s", err, userWithRole.Role.Role_type)
	return userWithRole, err // Return the loaded user and a nil error
}

// WhereUsername retrieves a user by username
func (repo *Repository) WhereUsername(username string, user *entity.User) error {
	err := repo.DB.Where("username = ?", username).First(user).Error
	return err
}

// WhereRoleID retrieves a role by user role ID
func (repo *Repository) WhereRoleID(Role_id string, role *entity.Role) error {

	err := repo.DB.Where("role_id = ?", Role_id).First(role).Error
	return err
}

func (repo *Repository) Preload_in_order() ([]entity.User, error) {
	var user_data []entity.User
	err := repo.DB.Preload("Role").Order("User_id asc").Find(&user_data).Error
	return user_data, err
}

func (repo *Repository) Role_in_Asc_order() ([]entity.Role, error) {
	var user_data []entity.Role
	err := repo.DB.Order("Role_id asc").Find(&user_data).Error
	return user_data, err
}

func (repo *Repository) Save(user *entity.User) error {
	return repo.DB.Save(user).Error
}

func (repo *Repository) Update(user *entity.User, update_user *entity.User) error {
	err := repo.DB.Model(user).Updates(update_user).Error
	return err
}

func (repo *Repository) Preload_Role_first(userWithRole *entity.User, User_id int) error {
	err := repo.DB.Preload("Role").First(userWithRole, User_id).Error
	return err
}

func (repo *Repository) Delete_user(user *entity.User) error {
	err := repo.DB.Delete(user).Error
	return err
}

func (repo *Repository) Delete_role(Role *entity.Role) error {
	err := repo.DB.Delete(Role).Error
	return err
}
