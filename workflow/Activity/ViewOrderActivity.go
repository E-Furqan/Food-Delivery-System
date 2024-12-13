package activity

import model "github.com/E-Furqan/Food-Delivery-System/Models"

func (act *Activity) ViewOrders(UserId uint, token string) (*[]model.UpdateOrder, error) {
	var userId model.UpdateOrder
	userId.DeliverDriverID = UserId
	Orders, err := act.OrderClient.ViewOrders(userId, token)
	if err != nil {
		return &[]model.UpdateOrder{}, err
	}
	return Orders, nil
}
