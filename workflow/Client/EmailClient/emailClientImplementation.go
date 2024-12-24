package EmailClient

import (
	"fmt"
	"net/smtp"

	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
)

func (email *EmailClient) SendEmail(orderID uint, orderStatus string, userEmail string) (string, error) {
	from := email.envVar.EmailAddressFrom
	password := email.envVar.EmailPassKey

	to := userEmail

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message, err := utils.GenerateEmail(orderID, orderStatus)
	if err != nil {
		return "", fmt.Errorf("failed to generate email: %w", err)
	}
	auth := smtp.PlainAuth("", from, password, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return "Email sent successfully", nil
}
