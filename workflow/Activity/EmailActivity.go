package activity

import (
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

func (act *Activity) FetchUserEmail(UserID uint, orderStatus string, token string) (string, error) {

	return "message", nil
}
