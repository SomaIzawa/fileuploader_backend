package controller

import (
	"file-uploader-api/awsmanager"
	"file-uploader-api/schema"
	"file-uploader-api/usecase"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IFileController interface {
	GetSignedURL(c echo.Context) error
	Download(c echo.Context) error
	DeleteFile(c echo.Context) error
}

type fileController struct {
	fu usecase.IFileUsecase
}

func NewFileController(fu usecase.IFileUsecase) IFileController {
	return &fileController{
		fu: fu,
	}
}

func (fc *fileController) GetSignedURL(c echo.Context) error {
	getSignedURLReq := schema.GetSignedURLReq{}
	if err := c.Bind(&getSignedURLReq); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	signedURL, err := awsmanager.GenerateSignedURL(getSignedURLReq.Url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, signedURL)
}

func (fc *fileController) Download(c echo.Context) error {
	id := c.Param("id")
	file, presignedURL, err := fc.fu.Download(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	filename := fmt.Sprintf("%s.%s", file.FileName, file.Type)
	return c.JSON(http.StatusOK, schema.FileDownloadRes{
		DownloadLink: presignedURL,
		FileName:     filename,
	})
}

func (fc *fileController) DeleteFile(c echo.Context) error {
	id := c.Param("id")
	if err := fc.fu.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusNoContent, nil)
}
