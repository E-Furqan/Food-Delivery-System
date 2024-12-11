package activity

import (
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

// func (act *Activity) Register(c *gin.Context) {

// 	var registrationData model.User

// 	if err := c.ShouldBindJSON(&registrationData); err != nil {
// 		utils.GenerateResponse(http.StatusBadRequest, c, "error", err.Error(), "", nil)
// 		return
// 	}

// 	if len(registrationData.Roles) > 0 && registrationData.ActiveRole == "" {
// 		var role model.Role
// 		if err := act.Repo.GetRole(registrationData.Roles[0].RoleId, &role); err != nil {
// 			utils.GenerateResponse(http.StatusInternalServerError, c, "error", "Role not found", "", nil)
// 			return
// 		}
// 		registrationData.ActiveRole = role.RoleType
// 		log.Print("active role set")
// 	}

// 	err := act.Repo.CreateUser(&registrationData)
// 	if err != nil {
// 		utils.GenerateResponse(http.StatusInternalServerError, c, "error", err.Error(), "", nil)
// 		return
// 	}

// 	c.JSON(http.StatusCreated, registrationData)
// }

func (act *Activity) RegisterCheckRole(registrationData model.User) (model.User, error) {
	if len(registrationData.Roles) > 0 && registrationData.ActiveRole == "" {
		var role model.Role
		if err := act.Repo.GetRole(registrationData.Roles[0].RoleId, &role); err != nil {
			// utils.GenerateResponse(http.StatusInternalServerError, c, "error", "Role not found", "", nil)
			return model.User{}, err
		}
		registrationData.ActiveRole = role.RoleType
		log.Print("active role set")
	}
	return registrationData, nil
}

func (act *Activity) CreateUser(registrationData model.User) (model.User, error) {
	err := act.Repo.CreateUser(&registrationData)
	if err != nil {
		log.Print("error userin activitiys", err)
		return model.User{}, err
	}

	return registrationData, nil
}
