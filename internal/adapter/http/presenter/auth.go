package presenter

type Token struct {
	AccessToken *string `json:"accessToken,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
}

func ToTokenResource(token *string) *Token {
	return &Token{
		AccessToken: token,
	}
}

type Authorize struct {
	Authorized bool   `json:"authorized"`
	JTI        string `json:"jti"`
	EXP        int64  `json:"exp"`
	ID         uint64 `json:"id"`
}
