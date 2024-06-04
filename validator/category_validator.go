package validator

import (
	"file-uploader-api/schema"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ICategoryValidator interface {
	CategoryCreateValidate(categoryReq schema.CreateCategoryReq) error
}

type categoryValidator struct{}

func NewCategoryValidator() ICategoryValidator {
	return &categoryValidator{}
}

func (cv *categoryValidator) CategoryCreateValidate(categoryReq schema.CreateCategoryReq) error {
	return validation.ValidateStruct(&categoryReq,
		validation.Field(
			&categoryReq.Name,
			validation.Required.Error("name is required"),
			validation.RuneLength(1, 30).Error("limited max 30 char"),
		),
	)
}
