package presenter

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
)

type User struct {
	ID        string  `json:"ID,omitempty" example:"8f4a1582-6a67-4d85-950b-2d17049c7385"`
	FirstName *string `json:"firstName,omitempty" example:"john"`
	LastName  *string `json:"lastName,omitempty" example:"doe"`
	Email     string  `json:"email,omitempty" example:"john.doe@gmail.com"`
	Status    string  `json:"status,omitempty" example:"ACTIVE"`
}

func prepareUser(user domain.User) User {
	return User{
		ID:        user.UUID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Status:    user.Status.String(),
	}
}

func ToUserResource(user *domain.User) *User {
	if user == nil {
		return nil
	}

	res := prepareUser(*user)
	return &res
}

func ToUserCollection(users []domain.User) []User {
	var result []User
	for _, user := range users {
		result = append(result, prepareUser(user))
	}

	return result
}
