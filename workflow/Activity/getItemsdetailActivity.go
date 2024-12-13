package activity

func (act *Activity) GetItems(orderID uint, orderStatus string) (string, error) {
	message, err := act.Email.EmailSender(orderID, orderStatus)
	if err != nil {
		return "", err
	}

	return message, nil
}
