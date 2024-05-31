package util

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
)

func AtoUint(s string) (uint, error) {
	convertedUInt64, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(convertedUInt64), nil
}

func GetFileMimeType(fileHeader *multipart.FileHeader) (string, error) {
	// ファイルを開く
	file, err := fileHeader.Open()
	if err != nil {
			return "", err
	}
	defer file.Close()

	// ファイルの先頭部分を読み取る
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
			return "", err
	}

	// MIME タイプを判別
	mimeType := http.DetectContentType(buffer)
	return mimeType, nil
}

func MimeTypeToType(fileHeader *multipart.FileHeader) (string, error) {
	mimeType, err := GetFileMimeType(fileHeader)

	if err != nil{
		return "", err
	}
	switch mimeType {
	case "image/jpeg":
		return "jpeg", nil
	case "image/png":
		return "png", nil
	default:
		return "", fmt.Errorf("unsupported mimetype")
	}
}
