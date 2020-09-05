package auth

import (
	"github.com/dgrijalva/jwt-go"
)

type AuthProfile struct {
	Email    string `json:"email" gorm:"column:email"`
	Password string `json:"password" gorm:"column:password"`
}

type JWTClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
}

type Response struct {
	httpStatusCode int
	Code           int         `json:"code"`
	Message        string      `json:"message"`
	Data           interface{} `json:"data"`
}

type ResponseBody struct {
	Token string `json:"token"`
}

func responseSuccess(httpStatusCode, code int, msg string, body interface{}) *Response {
	return &Response{httpStatusCode, code, msg, body}
}

func responseError(httpStatusCode, code int, msg string) *Response {
	return &Response{httpStatusCode, code, msg, nil}
}
