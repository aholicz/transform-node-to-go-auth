package auth

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var ErrDuplicatePrimaryKey = errors.New("Duplicate Primary Key")

type AuthDataSource struct {
	connection *gorm.DB
}

func NewAuthDataSource(conn *gorm.DB) *AuthDataSource {
	return &AuthDataSource{conn}
}

func (a *AuthDataSource) CreateAuthProfile(ctx context.Context, authProfile *AuthProfile) error {
	if err := a.connection.WithContext(ctx).Create(authProfile).Error; err != nil {
		errMessage := strings.ToLower(err.Error())
		if strings.Contains(errMessage, "duplicated") {
			return errors.Wrap(ErrDuplicatePrimaryKey, "email already exist in database")
		}
	}

	return nil
}

func (a *AuthDataSource) FindAuthProfile(ctx context.Context, email string) (*AuthProfile, error) {
	var authProfile *AuthProfile
	if err := a.connection.WithContext(ctx).Take(authProfile, email).Error; err != nil {
		return nil, errors.Wrap(err, "find auth profile from db")
	}

	return authProfile, nil
}
