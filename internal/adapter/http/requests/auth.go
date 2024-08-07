package requests

import "github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"

type AuthRegister struct {
	FirstName         *string `json:"firstName" binding:"required,regex_alpha,min=2,max=64" example:"john"`
	LastName          *string `json:"lastName" binding:"required,regex_alpha,min=2,max=64" example:"doe"`
	Email             string  `json:"email" binding:"required,email" example:"john.doe@gmail.com"`
	Password          string  `json:"password" binding:"required,min=8,max=64,password_complexity" example:"QWer123!@#"`
	ConfirmedPassword string  `json:"confirmedPassword" binding:"required,eqfield=Password" example:"QWer123!@#"`
}

func (r AuthRegister) ToUserDomain() domain.User {
	return domain.User{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Password:  &r.Password,
	}
}

type AuthLogin struct {
	Email    string `json:"email" binding:"required,email" example:"john.doe@gmail.com"`
	Password string `json:"password" binding:"required,password_complexity" example:"QWer123!@#"`
}

type AuthEmailOTPResend struct {
	Email string `json:"email" binding:"required,email" example:"john.doe@gmail.com"`
}

type AuthEmailOTPVerify struct {
	Email string `json:"email" binding:"required,email" example:"john.doe@gmail.com"`
	Token string `json:"token" binding:"required,token_length" example:"123456"`
}

type GoogleAuth struct {
	Email       string `json:"email" binding:"required,email" example:"john.doe@gmail.com"`
	AccessToken string `json:"accessToken" binding:"required" example:"123456789"`
}

type ForgetPassword struct {
	Email string `json:"email" binding:"required,email" example:"john@doe.com"`
}

type ResetPassword struct {
	Email             string `json:"email" binding:"required,email" example:"john@doe.com"`
	Token             string `json:"token" binding:"required,token_length" example:"123456"`
	Password          string `json:"password" binding:"required,password_complexity" example:"QWer123!@#"`
	ConfirmedPassword string `json:"confirmedPassword" binding:"required,eqfield=Password" example:"QWer123!@#"`
}

type AuthorizeRequest struct {
	RequiredPermissions []domain.PermissionKeyType `json:"requiredPermissions"`
}
