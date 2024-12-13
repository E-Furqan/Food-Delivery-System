package activity

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

func (act *Activity) UpdateOrderStatus(RestaurantId uint, order model.CombineOrderItem, token string) (string, error) {

	var updatedOrder model.UpdateOrder
	updatedOrder.OrderId = order.OrderId
	updatedOrder.OrderStatus = order.OrderStatus
	updatedOrder.UserID = order.UserID
	_, err := act.OrderClient.UpdateOrderStatus(updatedOrder, token)
	if err != nil {
		return "", err
	}

	return "s", nil
}

// func (act *Activity) UpdateOrderStatus(RestaurantId uint, order model.CombineOrderItem, token string) (string, error) {

// 	var updatedOrder model.UpdateOrder
// 	updatedOrder.OrderId = order.OrderId
// 	updatedOrder.OrderStatus = order.OrderStatus
// 	updatedOrder.UserID = order.UserID
// 	_, err := act.OrderClient.UpdateOrderStatus(updatedOrder, token)
// 	if err != nil {
// 		return "", err
// 	}

// 	return "s", nil
// }
