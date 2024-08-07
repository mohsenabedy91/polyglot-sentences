package requests

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"mime/multipart"
)

type CreateUserRequest struct {
	FirstName *string               `form:"firstName" binding:"required,regex_alpha,min=2,max=64" example:"John"`
	LastName  *string               `form:"lastName" binding:"required,regex_alpha,min=2,max=64" example:"Doe"`
	Email     string                `form:"email" binding:"required,email" example:"john.doe@gmail.com"`
	Avatar    *multipart.FileHeader `form:"avatar" binding:"required" swaggerignore:"true"`
}

func (r CreateUserRequest) ToUserDomain() domain.User {
	return domain.User{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
	}
}

type UserUUIDUri struct {
	UUIDStr string `uri:"userID" binding:"required,uuid" example:"8f4a1582-6a67-4d85-950b-2d17049c7385"`
}
