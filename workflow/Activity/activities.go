package activity

import (
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

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

// func (act *Activity) GetUser(UserID uint) (model.User, error) {
// 	var User model.User
// 	err := act.Repo.GetUser("user_id", UserID, &User)
// 	if err != nil {
// 		return model.User{}, err
// 	}
// 	return User, nil
// }

func (act *Activity) ViewOrders(UserId uint, token string) (*[]model.UpdateOrder, error) {
	var userId model.UpdateOrder
	userId.DeliverDriverID = UserId
	Orders, err := act.OrderClient.ViewOrders(userId, token)
	if err != nil {
		return &[]model.UpdateOrder{}, err
	}
	return Orders, nil
}
