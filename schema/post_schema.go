package schema

import (
	"file-uploader-api/model"
	"mime/multipart"
	"time"
)

type CreatePostReq struct {
	Title      string
	Comment    string
	CategoryID string
	UserID     string
	Thumnail   *multipart.FileHeader
	Files      []CreateFileReq
}

type PostRes struct {
	ID            uint        `json:"id"`
	Title         string      `json:"title"`
	Comment       string      `json:"comment"`
	ThumbnailType string      `json:"thumbnail_type"`
	User          UserRes     `json:"user"`
	Category      CategoryRes `json:"category"`
	Files         []FileRes   `json:"files"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

func PostModelToPostRes(post model.Post) PostRes {
	return PostRes{
		ID:            post.ID,
		Title:         post.Title,
		Comment:       post.Comment,
		ThumbnailType: post.ThumbnailType,
		User:          UserModelToUserRes(post.User),
		Category:      CategoryModelToCategoryRes(post.Category),
		Files:         FileModelListToFileResList(post.Files),
		CreatedAt:     post.CreatedAt,
		UpdatedAt:     post.UpdatedAt,
	}
}

func PostModelListToPostResList(posts []model.Post) []PostRes {
	postRes := make([]PostRes, len(posts))
	for i, post := range posts {
		postRes[i] = PostModelToPostRes(post)
	}
	return postRes
}
