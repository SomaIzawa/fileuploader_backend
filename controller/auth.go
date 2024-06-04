package controller

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetUserId(c echo.Context) (uint, error) {
	userObj := c.Get("user")
	if userObj == nil {
		return 0, errors.New("user token not found in context")
	}

	token, ok := userObj.(*jwt.Token)
	if !ok || token == nil {
		return 0, errors.New("user token is not of type *jwt.Token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("user claims type assertion failed")
	}

	userId, ok := claims["user_id"]
	if !ok {
		return 0, errors.New("user_id claim not found in token claims")
	}

	userIdFloat, ok := userId.(float64)
	if !ok {
		return 0, errors.New("user_id claim is not a float64")
	}
	return uint(userIdFloat), nil
}
