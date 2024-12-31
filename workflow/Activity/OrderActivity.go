package activityPac

import (
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

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

func (act *Activity) CalculateBill(CombineOrderItem model.CombineOrderItem, items []model.Items) (float64, error) {
	totalBill := 0.0

	for index, orderedItem := range CombineOrderItem.Items {
		totalBill += items[index].ItemPrice * float64(orderedItem.Quantity)
	}
	return totalBill, nil
}
