package port

type AuthService interface {
	GenerateToken(userUUIDStr string) (*string, error)
}
