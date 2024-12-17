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

func (act *Activity) CreateOrder(order model.CombineOrderItem, token string) (model.UpdateOrder, error) {

	Order, err := act.OrderClient.CreateOrder(order, token)
	if err != nil {
		log.Print("error from order activity: ", err)
		return model.UpdateOrder{}, err
	}
	log.Print("order from order activity:", Order)
	return Order, nil
}

func (act *Activity) CheckOrderStatus(orderID uint, token string) (string, error) {
	var OrderID model.OrderID
	OrderID.OrderID = orderID
	OrderStatus, err := act.OrderClient.FetchOrderStatus(OrderID, token)
	if err != nil {
		log.Print("error from check order activity: ", err)
		return "", err
	}
	log.Print("order status from check order activity:", OrderStatus)
	return OrderStatus, nil
}
