package jwt

type TokenPair struct {
	Refresh		string	`json:"-"`
	Access		string	`json:"access"`
}
