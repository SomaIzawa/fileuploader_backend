package controller

import (
	"crypto/x509"
	"encoding/pem"
	"file-uploader-api/schema"
	"file-uploader-api/usecase"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudfront/sign"
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
	signedURL, err := GenerateSignedURL(getSignedURLReq.Url)
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
		FileName: filename,
	})
}

func (fc *fileController) DeleteFile(c echo.Context) error {
	id := c.Param("id")
	if err := fc.fu.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusNoContent, nil)
}

func GenerateSignedURL(resourcePath string) (string, error) {
	// CloudFrontの設定
	keyPairID := os.Getenv("AWS_CLOUDFRONT_KEY_PAIR")

	// 秘密鍵の読み込み
	privateKeyBytes, err := os.ReadFile("private_key.pem")
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode([]byte(privateKeyBytes))
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	signer := sign.NewURLSigner(keyPairID, key)
	signedURL, err := signer.Sign(resourcePath, time.Now().Add(3*time.Second))
	if err != nil {
		return "", err
	}
	return signedURL, nil
}
