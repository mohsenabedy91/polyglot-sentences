package presenter

import (
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
)

type User struct {
	ID             string          `json:"id,omitempty" example:"8f4a1582-6a67-4d85-950b-2d17049c7385"`
	FirstName      *string         `json:"firstName,omitempty" example:"john"`
	LastName       *string         `json:"lastName,omitempty" example:"doe"`
	Email          string          `json:"email,omitempty" example:"john.doe@gmail.com"`
	Status         string          `json:"status,omitempty" example:"ACTIVE"`
	Gender         *string         `json:"gender,omitempty" example:"MALE"`
	AuthChallenges []AuthChallenge `json:"authChallenges"`
}

type AuthChallenge struct {
	Type   string `json:"type,omitempty" example:"Basic"`
	Status string `json:"status,omitempty" example:"ACTIVE"`
}

func PrepareUser(user *domain.User) *User {
	if user == nil || user.Base.UUID == uuid.Nil {
		return nil
	}

	return &User{
		ID:             user.Base.UUID.String(),
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		Status:         user.Status.String(),
		Gender:         helper.StringPtr(user.Gender.String()),
		AuthChallenges: PrepareAuthChallenge(user.AuthChallenges),
	}
}

func ToUserResource(user *domain.User) *User {
	return PrepareUser(user)
}

func ToUserCollection(users []*domain.User) []User {
	var response []User
	for _, userDetail := range users {
		result := PrepareUser(userDetail)
		if result != nil {
			response = append(response, *result)
		}
	}

	return response
}

func PrepareAuthChallenge(authChallenges []domain.AuthChallenge) []AuthChallenge {
	var response []AuthChallenge
	for _, authChallenge := range authChallenges {
		response = append(response, AuthChallenge{
			Type:   authChallenge.Type.String(),
			Status: authChallenge.Status.String(),
		})
	}

	return response
}

type TOTPKey struct {
	Secret string `json:"secret" example:"secret"`
	URL    string `json:"url" example:"otpauth_url"`
}

func ToTOTPResource(totp *domain.TOTPKey) *TOTPKey {
	if totp == nil {
		return nil
	}
	return &TOTPKey{
		Secret: totp.Secret,
		URL:    totp.URL,
	}
}
