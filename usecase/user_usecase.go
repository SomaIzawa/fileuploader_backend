package usecase

import (
	"file-uploader-api/model"
	"file-uploader-api/repository"
	"file-uploader-api/schema"
	"file-uploader-api/validator"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	SignUp(user schema.UserSignUpReq) (model.User, error)
	Login(user schema.UserLoginReq) (string, error)
}

type userUsecase struct {
	ur repository.IUserRepository
	uv validator.IUserValidator
}

func NewUserUsecase(ur repository.IUserRepository, uv validator.IUserValidator) IUserUsecase {
	return &userUsecase{ur, uv}
}

func (uu *userUsecase) SignUp(userReq schema.UserSignUpReq) (model.User, error) {
	if err := uu.uv.UserSignUpValidate(userReq); err != nil {
		return model.User{}, err
	}
	if userReq.Password != userReq.PasswordConfirm {
		return model.User{}, fmt.Errorf("password not match")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), 10)
	if err != nil {
		return model.User{}, err
	}
	newUser := model.User{
		Name:           userReq.Name,
		Email:          userReq.Email,
		HashedPassword: string(hash),
	}
	if err := uu.ur.CreateUser(&newUser); err != nil {
		return model.User{}, err
	}
	user := model.User{}
	if err := uu.ur.GetUserByEmail(&user, newUser.Email); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (uu *userUsecase) Login(userReq schema.UserLoginReq) (string, error) {
	if err := uu.uv.UserLoginValidate(userReq); err != nil {
		return "", err
	}
	storedUser := model.User{}
	if err := uu.ur.GetUserByEmail(&storedUser, userReq.Email); err != nil {
		return "", err
	}
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.HashedPassword), []byte(userReq.Password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   storedUser.ID,
		"user_name": storedUser.Name,
		"exp":       time.Now().Add(time.Hour * 12).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
