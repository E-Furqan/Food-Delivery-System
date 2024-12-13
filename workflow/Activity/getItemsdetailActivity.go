package activity

import (
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
)

func (act *Activity) GetItems(order model.CombineOrderItem, token string) (string, error) {
	var GetItem model.GetItems
	var itemCheck bool
	GetItem.RestaurantId = order.RestaurantId
	GetItem.ColumnName = "restaurant_id"
	GetItem.OrderType = "asc"

	items, err := act.ResClient.GetItems(GetItem)
	if err != nil {
		return "", err
	}

	itemCount := len(order.Items)
	for i := 0; i < itemCount; i++ {
		if !utils.ItemExists(order.Items[i], items) {
			itemCheck = true
		} else {
			itemCheck = false
			break
		}
	}

	if !itemCheck {
		act.SendEmail(order.OrderId, utils.Cancelled, token)
	}

	return "", nil
}
