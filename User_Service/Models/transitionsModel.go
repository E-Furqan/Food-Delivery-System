package model

func GetOrderTransitions() map[string]string {
	return map[string]string{
		"Waiting For Delivery Driver": "In for delivery",
		"In for delivery":             "Delivered",
		"Delivered":                   "Completed",
	}
}
