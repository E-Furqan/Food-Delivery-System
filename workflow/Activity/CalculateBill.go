package activity

import (
	"fmt"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

func (act *Activity) CalculateBill(CombineOrderItem model.CombineOrderItem, items []model.Items) (float64, error) {
	totalBill := 0.0

	for _, orderedItem := range CombineOrderItem.Items {
		var ItemPrice float64
		ItemFound := false

		for _, item := range items {
			if item.ItemId == orderedItem.ItemId {
				ItemPrice = item.ItemPrice
				ItemFound = true
				break
			}
		}

		if !ItemFound {
			continue
		}

		totalBill += ItemPrice * float64(orderedItem.Quantity)
	}
	if totalBill == 0 {
		return totalBill, fmt.Errorf("items are not of this restaurant")
	}
	return totalBill, nil
}
