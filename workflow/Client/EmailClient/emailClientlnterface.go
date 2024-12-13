package EmailClient

type EmailClient struct {
}

func NewClient() *EmailClient {
	return &EmailClient{}
}

type EmailClientInterface interface {
	EmailSender(orderID uint, orderStatus string) (string, error)
}
