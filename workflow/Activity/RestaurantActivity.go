package activity

import (
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
)

func (act *Activity) GetItems(order model.CombineOrderItem, token string) ([]model.Items, error) {

	items, err := act.ResClient.GetItems(order)
	if err != nil {
		log.Print("error from get items:", err)
		return []model.Items{}, err
	}
	return items, nil
}
