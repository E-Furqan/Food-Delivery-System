package EmailClient

import (
	"fmt"
	"log"
	"net/smtp"

	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
)

func (email *EmailClient) EmailSender(orderID uint, orderStatus string) (string, error) {
	from := "furqan.ali@emumba.com"
	password := "sqrf gefw qccw pqyr"

	to := "furqanali111366@gmail.com"

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

	log.Println("Email sent successfully!")
	return "Email sent successfully", nil
}
