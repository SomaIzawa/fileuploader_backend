package validator

import (
	"file-uploader-api/schema"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type IUserValidator interface {
	UserSignUpValidate(user schema.UserSignUpReq) error
	UserLoginValidate(user schema.UserLoginReq) error
}

type userValidator struct{}

func NewUserValidator() IUserValidator {
	return &userValidator{}
}

func (uv *userValidator) UserSignUpValidate(userReq schema.UserSignUpReq) error {
	return validation.ValidateStruct(&userReq,
		validation.Field(
			&userReq.Name,
			validation.Required.Error("name is required"),
			validation.RuneLength(1, 30).Error("limited max 30 char"),
		),
		validation.Field(
			&userReq.Email,
			validation.Required.Error("email is required"),
			validation.RuneLength(1, 30).Error("limited max 30 char"),
			is.Email.Error("is not valid email format"),
		),
		validation.Field(
			&userReq.Password,
			validation.Required.Error("password is required"),
			validation.RuneLength(6, 30).Error("limited min 6 max 30 char"),
		),
	)
}

func (uv *userValidator) UserLoginValidate(userReq schema.UserLoginReq) error {
	return validation.ValidateStruct(&userReq,
		validation.Field(
			&userReq.Email,
			validation.Required.Error("email is required"),
			validation.RuneLength(1, 30).Error("limited max 30 char"),
			is.Email.Error("is not valid email format"),
		),
		validation.Field(
			&userReq.Password,
			validation.Required.Error("password is required"),
			validation.RuneLength(6, 30).Error("limited min 6 max 30 char"),
		),
	)
}