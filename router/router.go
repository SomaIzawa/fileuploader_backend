package router

import (
	"file-uploader-api/awsmanager"
	"file-uploader-api/controller"
	"file-uploader-api/db"
	"file-uploader-api/repository"
	"file-uploader-api/usecase"
	"file-uploader-api/validator"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	echojwt "github.com/labstack/echo-jwt/v4"
)

func NewRouter() *echo.Echo {
	db := db.NewDB()
	am := awsmanager.NewAwsS3Manager()

	uv := validator.NewUserValidator()
	ur := repository.NewUserRepository(db)
	uu := usecase.NewUserUsecase(ur, uv)
	uc := controller.NewUserController(uu)

	cv := validator.NewCategoryValidator()
	cr := repository.NewCategoryRepository(db)
	cu := usecase.NewCategoryUsecase(cr, cv)
	cc := controller.NewCategoryController(cu)

	pv := validator.NewPostValidator()
	fv := validator.NewFileValidator()
	pr := repository.NewPostRepository(db)
	pu := usecase.NewPostUsecase(pr, ur, cr, pv, fv, am, db)
	pc := controller.NewPostController(pu)

	fr := repository.NewFileRepository(db)
	fu := usecase.NewFileUsecase(fr, am, db)
	fc := controller.NewFileController(fu)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowCredentials: true,
	}))
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     "/",
		CookieDomain:   os.Getenv("API_DOMAIN"),
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteNoneMode,
	}))
	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)
	e.GET("/csrf", uc.CsrfToken)

	c := e.Group("/categories")
	c.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
	}))
	c.POST("", cc.Create)
	c.GET("", cc.List)
	c.GET("/:categoryId", cc.GetById)

	p := e.Group("/posts")
	p.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
	}))
	p.POST("", pc.Create)
	p.GET("", pc.List)
	p.GET("/:postId", pc.GetById)

	f := e.Group("/files")
	f.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
	}))
	f.POST("/get-signed-url", fc.GetSignedURL)
	f.GET("/download/:id", fc.Download)
	f.DELETE("/:id", fc.DeleteFile)

	return e
}
