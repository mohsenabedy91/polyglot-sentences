package port

type AuthService interface {
	GenerateToken(userUUID string) (*string, error)
}
