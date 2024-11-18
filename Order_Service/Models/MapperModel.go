package model

func GetOrderTransitions() map[string]map[string]string {
	return map[string]map[string]string{
		"restaurant": {
			"order placed": "Accepted",
			"Accepted":     "In process",
			"In process":   "Waiting For Delivery Driver",
		},
		"delivery driver": {
			"Waiting For Delivery Driver": "In for delivery",
			"In for delivery":             "Delivered",
		},
		"user": {
			"Delivered": "Completed",
		},
	}
}
