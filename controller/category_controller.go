package controller

import (
	"file-uploader-api/schema"
	"file-uploader-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ICategoryController interface {
	List(c echo.Context) error
	GetById(c echo.Context) error
	Create(c echo.Context) error
}

type categoryController struct {
	cu usecase.ICategoryUsecase
}

func NewCategoryController(cu usecase.ICategoryUsecase) ICategoryController {
	return &categoryController{cu}
}

func (cc *categoryController) List(c echo.Context) error {
	categories, err := cc.cu.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, schema.CategoryModelListToCategoryResList(categories))
}

func (cc *categoryController) GetById(c echo.Context) error {
	id := c.Param("categoryId")
	category, err := cc.cu.Get(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, schema.CategoryModelToCategoryRes(category))
}

func (cc *categoryController) Create(c echo.Context) error {
	categoryCreateReq := schema.CreateCategoryReq{}
	if err := c.Bind(&categoryCreateReq); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	category, err := cc.cu.Create(categoryCreateReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, schema.CategoryModelToCategoryRes(category))
}
