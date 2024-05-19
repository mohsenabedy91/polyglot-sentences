package presenter

type Token struct {
	AccessToken *string `json:"accessToken,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
}

func ToTokenResource(token *string) *Token {
	return &Token{
		AccessToken: token,
	}
}
