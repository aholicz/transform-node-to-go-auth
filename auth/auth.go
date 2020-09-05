package auth

import (
	"context"
	"net/http"

	"gorm.io/gorm"

	"github.com/pkg/errors"
)

const (
	CodeSuccess             = 0
	CodeAuthProfileNotFound = 1000
	CodeInvalidPassword     = 1001
	CodeDuplicateEmail      = 1002
	CodeInternalError       = 9000
	CodeParseRequestError   = 9001

	DescSuccess             = "success"
	DescAuthProfileNotFound = "auth profile not found"
	DescInvalidPassword     = "invalid password"
	DescDuplicateEmail      = "this email already exist"
	DescInternalError       = "internal server error"
	DescParseRequestError   = "invalid request format"
)

var ErrInvalidPassword = errors.New(DescInvalidPassword)

type Service struct {
	datasource *AuthDataSource
}

func NewService(datasource *AuthDataSource) *Service {
	return &Service{datasource}
}

func (s *Service) CreateAuthProfile(ctx context.Context, authProfile *AuthProfile) (*Response, error) {
	if err := s.datasource.CreateAuthProfile(ctx, authProfile); err != nil {
		if errors.Cause(err) == ErrDuplicatePrimaryKey {
			return responseError(
				http.StatusUnauthorized,
				CodeDuplicateEmail,
				DescDuplicateEmail), err
		}
	}

	return responseSuccess(
		http.StatusCreated,
		CodeSuccess,
		DescSuccess, nil), nil
}

func (s *Service) Authenticate(ctx context.Context, request *AuthProfile) (*Response, error) {
	authProfile, err := s.datasource.FindAuthProfile(ctx, request.Email)
	if err != nil {
		if errors.Cause(err) == gorm.ErrRecordNotFound {
			return responseError(
				http.StatusUnauthorized,
				CodeAuthProfileNotFound,
				DescAuthProfileNotFound), err
		}

		return responseError(
			http.StatusUnauthorized,
			CodeInvalidPassword,
			DescInvalidPassword), err
	}

	hashPassword, err := hashPassword(authProfile.Password)
	if err != nil {
		return responseError(
			http.StatusInternalServerError,
			CodeInternalError,
			DescInternalError), ErrInvalidPassword
	}

	if valid := validatePassword(authProfile.Password, hashPassword); !valid {
		return responseError(
			http.StatusUnauthorized,
			CodeInvalidPassword,
			DescInvalidPassword), err
	}

	token, err := generateJwtToken(authProfile.Email)
	if err != nil {
		return responseError(
			http.StatusInternalServerError,
			CodeInternalError,
			DescInternalError), err
	}

	return responseSuccess(http.StatusOK, CodeSuccess, DescSuccess, &ResponseBody{token}), nil
}
