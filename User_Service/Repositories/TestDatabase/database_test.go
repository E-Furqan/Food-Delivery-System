package database_test

import (
	"fmt"
	"os"
	"testing"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	database "github.com/E-Furqan/Food-Delivery-System/Repositories"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println("Setup failed:", err)
		os.Exit(1)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	DB = ConfigTestDatabaseConnection()

	return m.Run(), nil
}

func setupRepository() *database.Repository {
	return &database.Repository{DB: DB}
}

func cleanupDB() {
	DB.Exec("DELETE FROM user_roles")
	DB.Exec("DELETE FROM roles")
	DB.Exec("DELETE FROM users")
}

func TestCreateRole(t *testing.T) {

	repo := setupRepository()
	cleanupDB()

	mockRole := model.Role{
		RoleId:   1,
		RoleType: "Admin",
	}
	err := repo.CreateRole(&mockRole)
	assert.NoError(t, err)

	var fetchRole model.Role

	err = DB.Where("role_id = ?", mockRole.RoleId).Find(&fetchRole).Error
	assert.NoError(t, err)

	assert.Equal(t, fetchRole.RoleType, mockRole.RoleType)

	err = repo.CreateRole(&mockRole)
	assert.Error(t, err)

	cleanupDB()
}

func TestCreateUser(t *testing.T) {

	repo := setupRepository()
	cleanupDB()

	mockUser := model.User{
		FullName:    "John Doe",
		Username:    "johndoe",
		Password:    "password",
		Email:       "johndoe@example.com",
		PhoneNumber: "1234567890",
		Address:     "123 Main St",
		RoleStatus:  "Active",
		ActiveRole:  "Admin",
	}
	err := repo.CreateUser(&mockUser)
	assert.NoError(t, err)

	var fetchUser model.User

	err = DB.Where("user_id = ?", mockUser.UserId).Find(&fetchUser).Error
	assert.NoError(t, err)

	assert.Equal(t, fetchUser.FullName, mockUser.FullName)
	assert.Equal(t, fetchUser.Username, mockUser.Username)
	assert.Equal(t, fetchUser.Password, mockUser.Password)
	assert.Equal(t, fetchUser.Email, mockUser.Email)
	assert.Equal(t, fetchUser.PhoneNumber, mockUser.PhoneNumber)
	assert.Equal(t, fetchUser.Address, mockUser.Address)
	assert.Equal(t, fetchUser.RoleStatus, mockUser.RoleStatus)
	assert.Equal(t, fetchUser.ActiveRole, mockUser.ActiveRole)

	mockUserSameEmail := model.User{
		FullName:    "John Doe1",
		Username:    "johndoe1",
		Password:    "password1",
		Email:       "johndoe@example.com",
		PhoneNumber: "12345678901",
		Address:     "123 Main St",
		RoleStatus:  "Active",
		ActiveRole:  "Admin",
	}

	err = repo.CreateUser(&mockUserSameEmail)
	assert.Error(t, err)

	mockUserSameUsername := model.User{
		FullName:    "John Doe12",
		Username:    "johndoe",
		Password:    "password12",
		Email:       "johndoe@example.com2",
		PhoneNumber: "123456789012",
		Address:     "123 Main St",
		RoleStatus:  "Active",
		ActiveRole:  "Admin",
	}

	err = repo.CreateUser(&mockUserSameUsername)
	assert.Error(t, err)

	mockUserSamePhoneNumber := model.User{
		FullName:    "John Doe121",
		Username:    "johndoe3",
		Password:    "password123",
		Email:       "johndoe@example.com23",
		PhoneNumber: "1234567890",
		Address:     "123 Main St",
		RoleStatus:  "Active",
		ActiveRole:  "Admin",
	}

	err = repo.CreateUser(&mockUserSamePhoneNumber)
	assert.Error(t, err)

	cleanupDB()
}

func TestLoadUserWithRoles(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	role1 := model.Role{RoleType: "Admin"}
	role2 := model.Role{RoleType: "User"}

	err := DB.Create(&role1).Error
	assert.NoError(t, err)

	err = DB.Create(&role2).Error
	assert.NoError(t, err)

	mockUser := model.User{
		FullName:    "John Doe",
		Username:    "johndoe",
		Password:    "password",
		Email:       "johndoe@example.com",
		PhoneNumber: "1234567890",
		Address:     "123 Main St",
		RoleStatus:  "Active",
		ActiveRole:  "Admin",
	}

	err = DB.Create(&mockUser).Error
	assert.NoError(t, err)

	userRole1 := model.UserRole{UserId: mockUser.UserId, RoleId: role1.RoleId}
	userRole2 := model.UserRole{UserId: mockUser.UserId, RoleId: role2.RoleId}

	err = DB.Create(&userRole1).Error
	assert.NoError(t, err)

	err = DB.Create(&userRole2).Error
	assert.NoError(t, err)

	var fetchedUser model.User
	fetchedUser.UserId = mockUser.UserId
	err = repo.LoadUserWithRoles(&fetchedUser)

	assert.NoError(t, err)
	assert.Equal(t, mockUser.FullName, fetchedUser.FullName)
	assert.Equal(t, mockUser.Username, fetchedUser.Username)
	assert.Equal(t, mockUser.Email, fetchedUser.Email)
	assert.Equal(t, 2, len(fetchedUser.Roles))
	assert.Contains(t, []string{role1.RoleType, role2.RoleType}, fetchedUser.Roles[0].RoleType)
	assert.Contains(t, []string{role1.RoleType, role2.RoleType}, fetchedUser.Roles[1].RoleType)

	cleanupDB()
}

func TestGetUser(t *testing.T) {

	repo := setupRepository()
	cleanupDB()

	mockUser := model.User{
		FullName:    "John Doe",
		Username:    "johndoe",
		Password:    "password",
		Email:       "johndoe@example.com",
		PhoneNumber: "1234567890",
		Address:     "123 Main St",
		RoleStatus:  "Active",
		ActiveRole:  "Admin",
	}
	err := DB.Create(&mockUser).Error
	assert.NoError(t, err)

	var fetchUser model.User

	err = repo.GetUser("user_id", mockUser.UserId, &fetchUser)
	assert.NoError(t, err)

	assert.Equal(t, fetchUser.FullName, mockUser.FullName)
	assert.Equal(t, fetchUser.Username, mockUser.Username)
	assert.Equal(t, fetchUser.Password, mockUser.Password)
	assert.Equal(t, fetchUser.Email, mockUser.Email)
	assert.Equal(t, fetchUser.PhoneNumber, mockUser.PhoneNumber)
	assert.Equal(t, fetchUser.Address, mockUser.Address)
	assert.Equal(t, fetchUser.RoleStatus, mockUser.RoleStatus)
	assert.Equal(t, fetchUser.ActiveRole, mockUser.ActiveRole)

	cleanupDB()
}

func TestGetRole(t *testing.T) {

	repo := setupRepository()
	cleanupDB()

	mockRole := model.Role{
		RoleId:   1,
		RoleType: "Admin",
	}
	err := DB.Create(&mockRole).Error
	assert.NoError(t, err)

	var fetchRole model.Role

	err = repo.GetRole(mockRole.RoleId, &fetchRole)
	assert.NoError(t, err)

	assert.Equal(t, fetchRole.RoleType, mockRole.RoleType)

	cleanupDB()
}

func TestFetchUsersWithRoles(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	role1 := model.Role{RoleType: "Admin"}
	role2 := model.Role{RoleType: "User"}

	err := DB.Create(&role1).Error
	assert.NoError(t, err)

	err = DB.Create(&role2).Error
	assert.NoError(t, err)

	mockUser1 := model.User{
		FullName:    "Alice Smith",
		Username:    "alice",
		Password:    "password123",
		Email:       "alice@example.com",
		PhoneNumber: "1234567890",
		Address:     "456 Elm St",
		RoleStatus:  "Active",
		ActiveRole:  "Admin",
	}

	mockUser2 := model.User{
		FullName:    "Bob Johnson",
		Username:    "bob",
		Password:    "password456",
		Email:       "bob@example.com",
		PhoneNumber: "0987654321",
		Address:     "789 Oak St",
		RoleStatus:  "Active",
		ActiveRole:  "Customer",
	}

	err = DB.Create(&mockUser1).Error
	assert.NoError(t, err)

	err = DB.Create(&mockUser2).Error
	assert.NoError(t, err)

	userRole1 := model.UserRole{UserId: mockUser1.UserId, RoleId: role1.RoleId}
	userRole2 := model.UserRole{UserId: mockUser2.UserId, RoleId: role2.RoleId}

	err = DB.Create(&userRole1).Error
	assert.NoError(t, err)

	err = DB.Create(&userRole2).Error
	assert.NoError(t, err)

	users, err := repo.FetchUsersWithRoles("full_name", "asc")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(users))

	assert.Contains(t, users[0].Roles, role1)
	assert.Contains(t, users[1].Roles, role2)

	cleanupDB()
}

func TestGetAllRoles(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	role1 := model.Role{RoleType: "Admin"}
	role2 := model.Role{RoleType: "User"}

	err := DB.Create(&role1).Error
	assert.NoError(t, err)

	err = DB.Create(&role2).Error
	assert.NoError(t, err)

	roles, err := repo.GetAllRoles()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(roles))
	assert.Contains(t, roles, role1)
	assert.Contains(t, roles, role2)

	cleanupDB()
}

func TestUpdateUser(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	mockOriginalUser := model.User{
		FullName:    "Jane Doe",
		Username:    "janedoe",
		Password:    "password123",
		Email:       "jane@example.com",
		PhoneNumber: "1234567890",
		Address:     "123 Main St",
		RoleStatus:  "Active",
		ActiveRole:  "User",
	}

	err := DB.Create(&mockOriginalUser).Error
	assert.NoError(t, err)

	mockUpdateUser := model.User{
		FullName:    "Jane Smith",
		Username:    "janesmith",
		Password:    "newpassword123",
		Email:       "jane.smith@example.com",
		PhoneNumber: "0987654321",
		Address:     "456 Oak St",
		RoleStatus:  "Inactive",
		ActiveRole:  "Admin",
	}

	err = repo.Update(&mockOriginalUser, &mockUpdateUser)
	assert.NoError(t, err)

	var updatedUser model.User
	err = DB.Preload("Roles").Where("user_id = ?", mockOriginalUser.UserId).First(&updatedUser).Error
	assert.NoError(t, err)

	assert.Equal(t, mockUpdateUser.FullName, updatedUser.FullName)
	assert.Equal(t, mockUpdateUser.Username, updatedUser.Username)
	assert.Equal(t, mockUpdateUser.Email, updatedUser.Email)
	assert.Equal(t, mockUpdateUser.PhoneNumber, updatedUser.PhoneNumber)
	assert.Equal(t, mockUpdateUser.Address, updatedUser.Address)
	assert.Equal(t, mockUpdateUser.RoleStatus, updatedUser.RoleStatus)
	assert.Equal(t, mockUpdateUser.ActiveRole, updatedUser.ActiveRole)

	cleanupDB()
}

func TestDeleteUser(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	user := model.User{
		FullName:    "Test User",
		Username:    "testuser",
		Password:    "password123",
		Email:       "test@example.com",
		PhoneNumber: "1234567890",
		Address:     "123 Test St",
		RoleStatus:  "Active",
		ActiveRole:  "User",
	}

	err := DB.Create(&user).Error
	assert.NoError(t, err)

	err = repo.DeleteUser(&user)
	assert.NoError(t, err)

	var deletedUser model.User
	err = repo.DB.Where("user_id = ?", user.UserId).First(&deletedUser).Error
	assert.Error(t, err)

	cleanupDB()
}

func TestDeleteRole(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	role := model.Role{
		RoleId:   1,
		RoleType: "User",
	}

	err := repo.DB.Create(&role).Error
	assert.NoError(t, err)

	err = repo.DeleteRole(&role)
	assert.NoError(t, err)

	var deletedRole model.Role
	err = repo.DB.Where("role_id = ?", role.RoleId).First(&deletedRole).Error
	assert.Error(t, err)

	cleanupDB()
}

func TestBulkCreateRoles(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	roles := []model.Role{
		{RoleType: "User"},
		{RoleType: "Admin"},
		{RoleType: "Manager"},
	}

	err := repo.BulkCreateRoles(roles)
	assert.NoError(t, err)

	var count int64
	err = repo.DB.Model(&model.Role{}).Count(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(len(roles)), count, "Expected number of roles to be created")

	for _, role := range roles {
		var fetchedRole model.Role
		err = repo.DB.Where("role_type = ?", role.RoleType).First(&fetchedRole).Error
		assert.NoError(t, err, "Role should exist in the database")
		assert.Equal(t, role.RoleType, fetchedRole.RoleType)
	}

	cleanupDB()
}

func TestDeleteUserRoleInfo(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	user := model.User{
		FullName:    "Test User",
		Username:    "testuser",
		Password:    "password123",
		Email:       "test@example.com",
		PhoneNumber: "1234567890",
		Address:     "123 Test St",
		RoleStatus:  "Active",
		ActiveRole:  "User",
	}
	err := DB.Create(&user).Error
	assert.NoError(t, err, "Error creating mock User")

	roles := []model.Role{
		{RoleId: 1, RoleType: "Admin"},
		{RoleId: 2, RoleType: "User"},
	}
	err = DB.Create(&roles).Error
	assert.NoError(t, err, "Error creating mock Roles")

	userRoles := []model.UserRole{
		{UserId: user.UserId, RoleId: roles[0].RoleId},
		{UserId: user.UserId, RoleId: roles[1].RoleId},
	}
	err = DB.Create(&userRoles).Error
	assert.NoError(t, err, "Error inserting mock data for UserRole")

	err = repo.DeleteUserRoleInfo(user.UserId, "user_user_id")
	assert.NoError(t, err)

	var remainingRoles []model.UserRole
	err = DB.Where("user_user_id = ?", user.UserId).Find(&remainingRoles).Error
	assert.NoError(t, err)
	assert.Empty(t, remainingRoles)

	userRoles1 := []model.UserRole{
		{UserId: user.UserId, RoleId: roles[0].RoleId},
		{UserId: user.UserId, RoleId: roles[1].RoleId},
	}
	err = DB.Create(&userRoles1).Error
	assert.NoError(t, err, "Error inserting mock data for UserRole")

	err = repo.DeleteUserRoleInfo(user.UserId, "role_role_id")
	assert.NoError(t, err)

	var remainingRoles1 []model.UserRole
	err = DB.Where("role_role_id = ?", user.UserId).Find(&remainingRoles).Error
	assert.NoError(t, err)
	assert.Empty(t, remainingRoles1)

	cleanupDB()
}

func TestAddUserRole(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	user := model.User{
		FullName:    "Test User",
		Username:    "testuser",
		Password:    "password123",
		Email:       "test@example.com",
		PhoneNumber: "1234567890",
		Address:     "123 Test St",
		RoleStatus:  "Active",
		ActiveRole:  "User",
	}
	err := DB.Create(&user).Error
	assert.NoError(t, err, "Error creating mock User")

	role := model.Role{RoleId: 1, RoleType: "Admin"}
	err = DB.Create(&role).Error
	assert.NoError(t, err, "Error creating mock Role")

	err = repo.AddUserRole(user.UserId, role.RoleId)
	assert.NoError(t, err, "Error adding new UserRole")

	var userRole model.UserRole
	err = DB.Where("user_user_id = ? AND role_role_id = ?", user.UserId, role.RoleId).First(&userRole).Error
	assert.NoError(t, err, "UserRole not found in the database")
	assert.Equal(t, user.UserId, userRole.UserId)
	assert.Equal(t, role.RoleId, userRole.RoleId)

	err = repo.AddUserRole(user.UserId, role.RoleId)
	assert.NoError(t, err, "Error adding duplicate UserRole")

	var userRoles []model.UserRole
	err = DB.Where("user_user_id = ? AND role_role_id = ?", user.UserId, role.RoleId).Find(&userRoles).Error
	assert.NoError(t, err)
	assert.Len(t, userRoles, 1, "Duplicate UserRole entry found")

	cleanupDB()
}

func TestUpdateUserActiveRole(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	user := model.User{
		FullName:    "Test User",
		Username:    "testuser",
		Password:    "password123",
		Email:       "test@example.com",
		PhoneNumber: "1234567890",
		Address:     "123 Test St",
		RoleStatus:  "Active",
		ActiveRole:  "User",
	}
	err := DB.Create(&user).Error
	assert.NoError(t, err, "Error creating mock User")

	user.ActiveRole = "Admin"
	err = repo.UpdateUserActiveRole(&user)
	assert.NoError(t, err, "Error updating User ActiveRole")

	var updatedUser model.User
	err = DB.Where("user_id = ?", user.UserId).First(&updatedUser).Error
	assert.NoError(t, err, "Error retrieving User after update")
	assert.Equal(t, "Admin", updatedUser.ActiveRole, "User ActiveRole was not updated correctly")

	cleanupDB()
}

func TestUpdateRoleStatus(t *testing.T) {
	repo := setupRepository()
	cleanupDB()

	user := model.User{
		FullName:    "Jane Doe",
		Username:    "janedoe",
		Password:    "password123",
		Email:       "janedoe@example.com",
		PhoneNumber: "555-555-5555",
		Address:     "789 Test Ave",
		RoleStatus:  "available",
		ActiveRole:  "User",
	}
	err := DB.Create(&user).Error
	assert.NoError(t, err, "Error creating mock User")

	user.RoleStatus = "not active"
	err = repo.UpdateRoleStatus(&user)
	assert.NoError(t, err, "Error updating role status")

	var updatedUser model.User
	err = DB.Where("user_id = ?", user.UserId).First(&updatedUser).Error
	assert.NoError(t, err, "Error retrieving updated user")
	assert.Equal(t, user.RoleStatus, updatedUser.RoleStatus)

	cleanupDB()
}

func TestConfigTestDatabaseConnection(t *testing.T) {

	db := ConfigTestDatabaseConnection()

	assert.NotNil(t, db)

	var userCount int64
	db.Model(&model.User{}).Count(&userCount)
	assert.Equal(t, int64(0), userCount, "Expected user count to be 0 after migrations")

	var roleCount int64
	db.Model(&model.Role{}).Count(&roleCount)
	assert.Equal(t, int64(0), roleCount, "Expected role count to be 0 after migrations")

	var userRoleCount int64
	db.Model(&model.UserRole{}).Count(&userRoleCount)
	assert.Equal(t, int64(0), userRoleCount, "Expected user role count to be 0 after migrations")
}
