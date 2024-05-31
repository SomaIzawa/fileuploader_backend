package schema

import (
	"file-uploader-api/model"
	"time"
)

type UserSignUpReq struct {
	Email string `json:"email"`
	Password string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type UserLoginReq struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserRes struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	AllowedAt time.Time `json:"allowed_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func UserModelToUserRes (user model.User) UserRes {
	return UserRes{
		user.ID,
		user.Name,
		user.Email,
		user.HashedPassword,
		user.AllowedAt,
		user.CreatedAt,
		user.UpdatedAt,
	}
}

