package usecase

import (
	"file-uploader-api/aws"
	"file-uploader-api/model"
	"file-uploader-api/repository"
	"file-uploader-api/schema"
	"file-uploader-api/util"
	"file-uploader-api/validator"
	"fmt"
	"os"

	"gorm.io/gorm"
)

type IPostUsecase interface {
	List() ([]model.Post, error)
	Get(id string) (model.Post, error)
	Create(req schema.CreatePostReq, userId uint)(model.Post, error)
}

type postUsecase struct {
	pr repository.IPostRepository
	ur repository.IUserRepository
	cr repository.ICategoryRepository
	pv validator.IPostValidator
	fv validator.IFileValidator
	db *gorm.DB
}

func NewPostUsecase(
	pr repository.IPostRepository, 
	ur repository.IUserRepository,
	cr repository.ICategoryRepository,
	pv validator.IPostValidator, 
	fv validator.IFileValidator, 
	db *gorm.DB) IPostUsecase {
	return &postUsecase{
		pr: pr,
		ur: ur,
		cr: cr,
		pv: pv,
		fv: fv,
		db: db,
	}
}

func (pu *postUsecase) List() ([]model.Post, error) {
	posts := []model.Post{}
	if err := pu.pr.GetPosts(&posts); err != nil {
		return []model.Post{}, err
	}
	return posts, nil
}

func (pu *postUsecase) Get(id string) (model.Post, error) {
	post := model.Post{}
	uintId, err := util.AtoUint(id);
	if err != nil {
		return model.Post{}, err
	}
	if err != pu.pr.GetPost(&post, uintId) {
		return model.Post{}, err
	}
	return post, nil
}

func (pu *postUsecase) Create(req schema.CreatePostReq, userId uint) (model.Post, error) {
	// awsの準備
	s := aws.NewS3Session()
	uploader := aws.CreateUploader(s)
	awsBucketName := os.Getenv("AWS_BUCKET_NAME")

	//バリデーション

	// バリデーターを用いたバリデーション
	if err := pu.pv.PostCreateValidate(req); err != nil {
		return model.Post{}, err
	}
	for _, fileReq := range req.Files {
		if err := pu.fv.FileCreateValidate(fileReq); err != nil {
			return model.Post{}, err
		}	
	}
	//カテゴリーIDの有効性
	categoryId, err := util.AtoUint(req.CategoryID)
	if err != nil {
		return model.Post{}, err
	}
	if err := pu.cr.GetCategory(&model.Category{}, categoryId); err != nil {
		return model.Post{}, err
	}
	//ユーザーIDの有効性
	if err := pu.ur.GetUser(&model.User{}, userId); err != nil {
		return model.Post{}, err
	}

	// 新規追加するデータの用意
	thumbnailFileType, err := util.MimeTypeToType(req.Thumnail)
	if err != nil {
		return model.Post{}, err
	}
	uintCategoryID, err := util.AtoUint(req.CategoryID)
	if err != nil {
		return model.Post{}, err
	}
	newPost := model.Post{
		Title: req.Title,
		Comment: req.Comment,
		ThumbnailType: thumbnailFileType,
		UserId: userId,
		CategoryId: uintCategoryID,
		Files: func() []model.File {
			files := make([]model.File, len(req.Files))
			for i, file := range req.Files {
				fileType, _ := util.MimeTypeToType(file.File)
				files[i] = model.File{
					FileName: file.Name,
					Type: fileType,
					UserId: userId,
				}
			}
			return files
		}(),
	}

	tx := pu.db.Begin()

	// 投稿情報をDBに保存
	if err := pu.pr.Create(&newPost); err != nil {
		return model.Post{}, nil
	}

	// アップロード
	if err := aws.S3Upload(uploader, awsBucketName, fmt.Sprintf("t%d.%s", newPost.ID, newPost.ThumbnailType), *req.Thumnail); err != nil {
		tx.Rollback()
		return model.Post{}, err	
	} 
	for i, file := range req.Files {
		if err := aws.S3Upload(uploader, awsBucketName, fmt.Sprintf("%d.%s", newPost.Files[i].ID, newPost.Files[i].Type), *file.File); err != nil {
			tx.Rollback()
			return model.Post{}, err	
		} 
	}

	// トランザクションをコミット
	if err := tx.Commit().Error; err != nil {
		return model.Post{}, err
	}
	
	return newPost, nil
}