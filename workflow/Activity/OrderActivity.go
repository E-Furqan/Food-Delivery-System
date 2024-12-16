package activity

import (
	"log"

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

	return "order status updated", nil
}

func (act *Activity) ViewOrders(UserId uint, token string) (*[]model.UpdateOrder, error) {
	var userId model.UpdateOrder
	userId.DeliverDriverID = UserId
	Orders, err := act.OrderClient.ViewOrders(userId, token)
	if err != nil {
		return &[]model.UpdateOrder{}, err
	}
	return Orders, nil
}

func (act *Activity) CreateOrder(order model.CombineOrderItem, token string) (uint, error) {

	OrderID, err := act.OrderClient.CreateOrder(order, token)
	if err != nil {
		log.Print("error from order activity: ", err)
		return 0, err
	}
	return OrderID.OrderId, nil
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
