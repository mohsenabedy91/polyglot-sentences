package presenter

import (
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
)

type User struct {
	ID        string  `json:"id,omitempty" example:"8f4a1582-6a67-4d85-950b-2d17049c7385"`
	FirstName *string `json:"firstName,omitempty" example:"john"`
	LastName  *string `json:"lastName,omitempty" example:"doe"`
	Email     string  `json:"email,omitempty" example:"john.doe@gmail.com"`
	Status    string  `json:"status,omitempty" example:"ACTIVE"`
}

func PrepareUser(user *domain.User) *User {
	if user == nil || user.Base.UUID == uuid.Nil {
		return nil
	}

	return &User{
		ID:        user.Base.UUID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Status:    user.Status.String(),
	}
}

func ToUserResource(user *domain.User) *User {
	return PrepareUser(user)
}

func ToUserCollection(users []*domain.User) []User {
	var response []User
	for _, user := range users {
		result := PrepareUser(user)
		if result != nil {
			response = append(response, *result)
		}
	}

	return response
}
