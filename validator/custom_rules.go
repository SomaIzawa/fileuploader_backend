package validator

import (
	"errors"
	"file-uploader-api/util"
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func FileType(acceptTypes []string) validation.RuleFunc {
	return func(value interface{}) error {
		file, ok := value.(*multipart.FileHeader)
		if !ok {
			return errors.New("invalid file type")
		}

		fileType, err := util.GetFileMimeType(file)
		if err != nil {
			return errors.New("failed to get file type by *multipart.FileHeader")
		}

		for _, acceptType := range acceptTypes {
			if fileType == acceptType {
				return nil
			}
		}
		return errors.New("unsupported file type")
	}
}