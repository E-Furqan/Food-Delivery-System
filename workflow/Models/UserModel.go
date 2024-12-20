package model

type UserClientEnv struct {
	Fetch_email_URL string
	BASE_URL        string
	USER_PORT       string
}

type UserEmail struct {
	Email string `json:"email"`
}
