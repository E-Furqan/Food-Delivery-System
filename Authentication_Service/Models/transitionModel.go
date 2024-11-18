package model

func GetOrderTransitions() map[string]string {
	return map[string]string{
		"order placed": "Accepted",
		"Accepted":     "In process",
		"In process":   "Waiting For Delivery Driver",
	}
}
