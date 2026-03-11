package auth

type CreateMagicLinkPayload struct {
	Email	string	`json:"email"`
}

type CreateMagicLinkResponse struct {
	Token	string `json:"token"`
}

