package activity

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

func (act *Activity) OrderConfirmationEmail(orderID uint, orderStatus string) (string, error) {
	from := "furqan.ali@emumba.com"
	password := "sqrf gefw qccw pqyr"

	to := "furqanali111366@gmail.com"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	var subject, body string

	switch strings.ToLower(orderStatus) {
	case "cancelled":
		subject = "Subject: Order Cancellation\n"
		body = fmt.Sprintf("We regret to inform you that your order with Order ID %v has been cancelled. If you have any questions, please contact support.", orderID)
	case "accepted":
		subject = "Subject: Order Confirmation\n"
		body = fmt.Sprintf("Great news! Your order with Order ID %v has been confirmed. Thank you for choosing us.", orderID)
	case "completed":
		subject = "Subject: Order Completed\n"
		body = fmt.Sprintf("Congratulations! Your order with Order ID %v has been successfully completed. We would appreciate it if you could leave a review. Thank you!", orderID)
	default:
		log.Println("Invalid order status provided.")
		return "", fmt.Errorf("invalid order status: %s", orderStatus)
	}

	message := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	log.Println("Email sent successfully!")
	return "Email sent successfully", nil
}
