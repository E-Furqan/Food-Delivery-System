package activity

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

func (act *Activity) CalculateBill(CombineOrderItem model.CombineOrderItem, items []model.Items) (float64, error) {
	totalBill := 0.0

	for index, orderedItem := range CombineOrderItem.Items {
		totalBill += items[index].ItemPrice * float64(orderedItem.Quantity)
	}
	return totalBill, nil
}
