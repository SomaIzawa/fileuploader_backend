package schema

import (
	"file-uploader-api/model"
	"mime/multipart"
	"time"
)

type CreateFileReq struct {
	File   *multipart.FileHeader
	Name   string
	UserId string
}

type GetSignedURLReq struct {
	Url string `json:"url"`
}

type FileRes struct {
	ID        uint      `json:"id"`
	FileName  string    `json:"file_name"`
	Type      string    `json:"type"`
	User      UserRes   `json:"user"`
	Post      PostRes   `json:"post"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FileDownloadRes struct {
	DownloadLink string `json:"download_link"`
	FileName     string `json:"file_name"`
}

func FileModelToFileRes(file model.File) FileRes {
	return FileRes{
		ID:        file.ID,
		FileName:  file.FileName,
		Type:      file.Type,
		User:      UserModelToUserRes(file.User),
		Post:      PostModelToPostRes(file.Post),
		CreatedAt: file.CreatedAt,
		UpdatedAt: file.UpdatedAt,
	}
}

func FileModelListToFileResList(files []model.File) []FileRes {
	fileRes := make([]FileRes, len(files))
	for i, file := range files {
		fileRes[i] = FileModelToFileRes(file)
	}
	return fileRes
}
