package requests

import "github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"

type AuthRegister struct {
	FirstName         *string `json:"firstName" binding:"required,regex_alpha,min=2,max=64" example:"john"`
	LastName          *string `json:"lastName" binding:"required,regex_alpha,min=2,max=64" example:"doe"`
	Email             string  `json:"email" binding:"required,email" example:"john.doe@gmail.com"`
	Password          string  `json:"password" binding:"required,min=8,max=64,password_complexity" example:"QWer123!@#"`
	ConfirmedPassword string  `json:"confirmedPassword" binding:"required,eqfield=Password" example:"QWer123!@#"`
	Gender            string  `json:"gender" binding:"required,oneof=MALE FEMALE OTHER PREFER_NOT_TO_SAY" example:"PREFER_NOT_TO_SAY"`
}

func (r AuthRegister) ToUserDomain() domain.User {
	return domain.User{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Password:  &r.Password,
		Gender:    domain.ToUserGenderType(r.Gender),
	}
}

type AuthLogin struct {
	Email    string `json:"email" binding:"required,email" example:"john.doe@gmail.com"`
	Password string `json:"password" binding:"required,password_complexity" example:"QWer123!@#"`
}

type AuthChallenge struct {
	Type string `json:"type" binding:"required" example:"totp"`
	Code string `json:"code" binding:"required" example:"123456"`
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
