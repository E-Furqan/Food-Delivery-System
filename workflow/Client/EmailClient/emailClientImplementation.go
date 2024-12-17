package EmailClient

import (
	"fmt"
	"net/smtp"

	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
)

func (email *EmailClient) EmailSender(orderID uint, orderStatus string) (string, error) {
	from := utils.FurqanEmail
	password := utils.FurqanEmailPassKey

	to := utils.SenderEmail

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message, err := utils.EmailGenerator(orderID, orderStatus)
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
