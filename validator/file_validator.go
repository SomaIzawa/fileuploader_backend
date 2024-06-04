package validator

import (
	"file-uploader-api/schema"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IFileValidator interface {
	FileCreateValidate(fileReq schema.CreateFileReq) error
}

type fileValidator struct{}

func NewFileValidator() IFileValidator {
	return &fileValidator{}
}

var acceptTypes = []string{"image/jpeg", "image/png"}

func (cv *fileValidator) FileCreateValidate(fileReq schema.CreateFileReq) error {
	return validation.ValidateStruct(&fileReq,
		validation.Field(
			&fileReq.Name,
			validation.Required.Error("name is required"),
			validation.RuneLength(1, 30).Error("limited max 30 char"),
		),
		validation.Field(
			&fileReq.File,
			validation.By(FileType(acceptTypes)),
		),
	)
}
