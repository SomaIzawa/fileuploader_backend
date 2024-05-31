package validator

import (
	"file-uploader-api/schema"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IPostValidator interface {
	PostCreateValidate(postReq schema.CreatePostReq) error
}

type postValidator struct{}

func NewPostValidator() IPostValidator {
	return &postValidator{}
}

var thumbnailAcceptTypes = []string{"image/jpeg", "image/png"}

func (pv *postValidator) PostCreateValidate(postReq schema.CreatePostReq) error {
	return validation.ValidateStruct(&postReq,
		validation.Field(
			&postReq.Title,
			validation.Required.Error("title is required"),
			validation.RuneLength(1, 30).Error("limited max 30 char"),
		),
		validation.Field(
			&postReq.Comment,
			validation.Required.Error("comment is required"),
			validation.RuneLength(1, 300).Error("limited max 30 char"),
		),
		validation.Field(
			&postReq.Thumnail,
			validation.By(FileType(thumbnailAcceptTypes)),
		),
	)
}