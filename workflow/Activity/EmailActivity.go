package activity

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
)

func (act *Activity) SendEmail(orderID uint, orderStatus string, token string, userEmail string) (string, error) {

	if orderStatus == utils.Cancelled {
		updatedOrder := utils.UpdateOrderStatusTOCancel(orderID)
		act.OrderClient.UpdateOrderStatus(updatedOrder, token)

	}

	message, err := act.Email.EmailSender(orderID, orderStatus, userEmail)
	if err != nil {
		return "", err
	}

	return message, nil
}

func (act *Activity) FetchUserEmail(token string) (*model.UserEmail, error) {

	email, err := act.UserClient.FetchEmail(token)
	if err != nil {
		return &model.UserEmail{}, err
	}

	return email, nil
}
