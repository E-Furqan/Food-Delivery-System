package database

import (
	"fmt"
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

func (repo *Repository) CreateUser(user *model.User) error {
	result := repo.DB.Create(user)
	repo.LoadUserWithRoles(user)
	return result.Error
}

func (repo *Repository) LoadUserWithRoles(user *model.User) error {

	err := repo.DB.Preload("Roles").First(user, user.UserId).Error
	if err != nil {
		log.Printf("Error loading user with roles: %v", err)
		return err
	}

	return nil
}

func (repo *Repository) GetRole(RoleId uint, role *model.Role) error {

	err := repo.DB.Where("role_id = ?", RoleId).First(role).Error
	return err
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
