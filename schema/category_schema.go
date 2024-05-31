package schema

import (
	"file-uploader-api/model"
	"time"
)

type CreateCategoryReq struct {
	Name string `json:"name"`	
}

type CategoryRes struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CategoryModelToCategoryRes (category model.Category) CategoryRes {
	return CategoryRes{
		category.ID,
		category.Name,
		category.CreatedAt,
		category.UpdatedAt,
	}
}

func CategoryModelListToCategoryResList (categories []model.Category) []CategoryRes {
	list := make([]CategoryRes, len(categories))
	for i, category := range categories {
		list[i] = CategoryModelToCategoryRes(category)
	}
	return list
}