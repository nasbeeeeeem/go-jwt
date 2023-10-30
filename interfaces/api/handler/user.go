package handler

import (
	"net/http"

	"github.com/FarStep131/go-jwt/myerror"
	"github.com/FarStep131/go-jwt/usecase"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	HandlerSing(c echo.Context)
	HandlerLogin(c echo.Context)
	HandlerLogout(c echo.Context)
}

type handler struct {
	useCase usecase.UseCase
}

// HandlerSing implements Handler.
func (h *handler) HandlerSing(c echo.Context) {
	type (
		request struct {
			Username string `json:"username" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=8"`
		}
		response struct {
			ID       int64  `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		}
	)

	requestBody := new(request)
	if err := c.Bind(&requestBody); err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}

	user, err := h.useCase.Singup(c.Request().Context(), requestBody.Username, requestBody.Email, requestBody.Password)
	if err != nil {
		switch e := err.(type) {
		case *myerror.InternalSeverError:
			c.JSON(http.StatusInternalServerError, e.Err.Error())
			return
		case *myerror.BadRequestError:
			c.JSON(http.StatusBadRequest, e.Err.Error())
			return
		default:
			c.JSON(http.StatusInternalServerError, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, &response{
		ID:       int64(user.ID),
		Username: user.UserName,
		Email:    user.Email,
	})
}

// HandlerLogin implements Handler.
func (h *handler) HandlerLogin(c echo.Context) {
	type (
		request struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}
		response struct {
			ID       int64  `json:"id"`
			Username string `json:"username"`
		}
	)

	requestBody := new(request)

	if err := c.Bind(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	singedString, user, err := h.useCase.Login(c.Request().Context(), requestBody.Email, requestBody.Password)
	if err != nil {
		switch e := err.(type) {
		case *myerror.InternalSeverError:
			c.JSON(http.StatusInternalServerError, e.Err.Error())
			return
		case *myerror.BadRequestError:
			c.JSON(http.StatusBadRequest, e.Err.Error())
			return
		default:
			c.JSON(http.StatusInternalServerError, e.Error())
		}
	}

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    singedString,
		MaxAge:   60 * 60 * 24, // 24h
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	}
	c.SetCookie(cookie)

	c.JSON(http.StatusOK, &response{
		ID:       int64(user.ID),
		Username: user.UserName,
	})
}

// HandlerLogout implements Handler.
func (*handler) HandlerLogout(c echo.Context) {
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		MaxAge:   -1,
		Path:     "",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	}
	c.SetCookie(cookie)
	c.JSON(http.StatusOK, "logout successful")
}

func NewHandler(usecase usecase.UseCase) Handler {
	return &handler{
		useCase: usecase,
	}
}
