package requests

import "github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"

type AuthRegister struct {
	FirstName         string `json:"firstName" binding:"required,regex_alpha,min=2,max=64" example:"john"`
	LastName          string `json:"lastName" binding:"required,regex_alpha,min=2,max=64" example:"doe"`
	Email             string `json:"email" binding:"required,email" example:"john.doe@gmail.com"`
	Password          string `json:"password" binding:"required,min=8,max=64,password_complexity" example:"QWer123!@#"`
	ConfirmedPassword string `json:"confirmedPassword" binding:"required,eqfield=Password" example:"QWer123!@#"`
}

func (r AuthRegister) ToDomain() domain.User {
	return domain.User{
		FirstName: &r.FirstName,
		LastName:  &r.FirstName,
		Email:     &r.Email,
		Password:  &r.Password,
	}
}

type AuthLogin struct {
	Email    string `json:"email" binding:"required,email" example:"john.doe@gmail.com"`
	Password string `json:"password" binding:"required,password_complexity" example:"QWer123!@#"`
}
