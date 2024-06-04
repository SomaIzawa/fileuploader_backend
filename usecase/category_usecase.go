package usecase

import (
	"file-uploader-api/model"
	"file-uploader-api/repository"
	"file-uploader-api/schema"
	"file-uploader-api/util"
	"file-uploader-api/validator"
)

type ICategoryUsecase interface {
	List() ([]model.Category, error)
	Get(id string) (model.Category, error)
	Create(req schema.CreateCategoryReq) (model.Category, error)
}

type categoryUsecase struct {
	cr repository.ICategoryRepository
	cv validator.ICategoryValidator
}

func NewCategoryUsecase(ur repository.ICategoryRepository, uv validator.ICategoryValidator) ICategoryUsecase {
	return &categoryUsecase{ur, uv}
}

func (cu *categoryUsecase) List() ([]model.Category, error) {
	categories := []model.Category{}
	if err := cu.cr.GetCategories(&categories); err != nil {
		return []model.Category{}, err
	}
	return categories, nil
}

func (cu *categoryUsecase) Get(id string) (model.Category, error) {
	category := model.Category{}
	uintId, err := util.AtoUint(id)
	if err != nil {
		return model.Category{}, err
	}
	if err != cu.cr.GetCategory(&category, uintId) {
		return model.Category{}, err
	}
	return category, nil
}

func (cu *categoryUsecase) Create(req schema.CreateCategoryReq) (model.Category, error) {
	if err := cu.cv.CategoryCreateValidate(req); err != nil {
		return model.Category{}, err
	}
	newCategory := model.Category{Name: req.Name}
	id, err := cu.cr.CreateCategory(&newCategory)
	if err != nil {
		return model.Category{}, err
	}
	category := model.Category{}
	if err := cu.cr.GetCategory(&category, id); err != nil {
		return model.Category{}, err
	}
	return category, nil
}
