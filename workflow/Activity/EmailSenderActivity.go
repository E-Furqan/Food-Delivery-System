package activity

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
)

func (act *Activity) SendEmail(orderID uint, orderStatus string, token string) (string, error) {
	message, err := act.Email.EmailSender(orderID, orderStatus)
	if err != nil {
		return "", err
	}

	if orderStatus == utils.Cancelled {
		var order model.CombineOrderItem
		order.OrderId = orderID
		order.OrderStatus = orderStatus
		order.UserID = 0
		updatedOrder := utils.UpdateOrderStatusTOCancel(order)
		act.OrderClient.UpdateOrderStatus(updatedOrder, token)
	}
	return message, nil
}
