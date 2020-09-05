package auth

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/tOnkowzl/libs/logx"
)

type Handler struct {
	service *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc}
}

func (h *Handler) CreateAuthProfile(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		authProfile *AuthProfile
	)

	if err := c.Bind(authProfile); err != nil {
		logx.WithContext(ctx).Errorf("binding request %+v", errors.WithStack(err))
		return c.JSON(http.StatusBadRequest, responseError(
			http.StatusBadRequest,
			CodeParseRequestError,
			DescParseRequestError))
	}

	if missing, msg := isMissingRequired(authProfile); missing {
		return c.JSON(http.StatusBadRequest, responseError(
			http.StatusBadRequest,
			CodeMissingRequiredFields,
			msg))
	}

	resp, err := h.service.CreateAuthProfile(ctx, authProfile)
	if err != nil {
		logx.WithContext(ctx).Errorf("%+v", err)
	}

	return c.JSON(resp.httpStatusCode, resp)
}

func (h *Handler) Authenticate(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		authProfile *AuthProfile
	)

	if err := c.Bind(authProfile); err != nil {
		logx.WithContext(ctx).Errorf("binding request %+v", errors.WithStack(err))
		return c.JSON(http.StatusBadRequest, responseError(
			http.StatusBadRequest,
			CodeParseRequestError,
			DescParseRequestError))
	}

	if missing, msg := isMissingRequired(authProfile); missing {
		return c.JSON(http.StatusBadRequest, responseError(
			http.StatusBadRequest,
			CodeMissingRequiredFields,
			msg))
	}

	resp, err := h.service.Authenticate(ctx, authProfile)
	if err != nil {
		logx.WithContext(ctx).Errorf("%+v", err)
	}

	return c.JSON(resp.httpStatusCode, resp)
}
