package controller

import (
	"file-uploader-api/schema"
	"file-uploader-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IPostController interface {
	List(c echo.Context) error
	GetById(c echo.Context) error
	Create(c echo.Context) error
}

type postController struct {
	pu usecase.IPostUsecase
}

func NewPostController(pu usecase.IPostUsecase) IPostController {
	return &postController{
		pu: pu,
	}
}

func (pc *postController) List(c echo.Context) error {
	posts, err := pc.pu.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, schema.PostModelListToPostResList(posts))
}

func (pc *postController) GetById(c echo.Context) error {
	id := c.Param("postId")
	post, err := pc.pu.Get(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, schema.PostModelToPostRes(post))
}

func (pc *postController) Create(c echo.Context) error {
	// リクエスト解析
	var createPostReq schema.CreatePostReq

	createPostReq.Title = c.FormValue("title")
	createPostReq.Comment = c.FormValue("comment")
	createPostReq.CategoryID = c.FormValue("category_id")

	// サムネイル
	thumbnail, err := c.FormFile("thumbnail")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	createPostReq.Thumnail = thumbnail

	// ファイル
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]
	fileNames := c.Request().Form["filenames"]

	createFileReqs := make([]schema.CreateFileReq, len(files))

	for i, _ := range createFileReqs {
		createFileReqs[i] = schema.CreateFileReq{
			File: files[i],
			Name: fileNames[i],
		}
	}

	createPostReq.Files = createFileReqs

	// ログインユーザの取得
	userId, err := GetUserId(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, err.Error())
	}

	// ユースケース実行
	createdPost, err := pc.pu.Create(createPostReq, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, createdPost)
}
