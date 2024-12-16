package activity

import (
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
)

func (act *Activity) GetItems(order model.CombineOrderItem, token string) ([]model.Items, error) {

	items, err := act.ResClient.GetItems(order)
	if err != nil {
		log.Print("error from get items:", err)
		act.SendEmail(order.OrderId, utils.Cancelled, token)
		return []model.Items{}, err
	}
	log.Print("items;  ", items)

	return items, nil
}
