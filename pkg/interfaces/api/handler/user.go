package handler

import (
	"net/http"

	"go-jwt/pkg/myerror"
	"go-jwt/pkg/usecase"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	HandlerSing(c *gin.Context)
	HandlerLogin(c *gin.Context)
	HandlerLogout(c *gin.Context)
}

type handler struct {
	useCase usecase.UseCase
}

func NewHandler(usecase usecase.UseCase) Handler {
	return &handler{
		useCase: usecase,
	}
}

// HandlerSing implements Handler.
func (h *handler) HandlerSing(c *gin.Context) {
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

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.useCase.Singup(c.Request.Context(), requestBody.Username, requestBody.Email, requestBody.Password)
	if err != nil {
		switch e := err.(type) {
		case *myerror.InternalSeverError:
			c.JSON(http.StatusInternalServerError, gin.H{"error": e.Err.Error()})
			return
		case *myerror.BadRequestError:
			c.JSON(http.StatusBadRequest, gin.H{"error": e.Err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
func (h *handler) HandlerLogin(c *gin.Context) {
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

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	singedString, user, err := h.useCase.Login(c.Request.Context(), requestBody.Email, requestBody.Password)

	if err != nil {
		switch e := err.(type) {
		case *myerror.InternalSeverError:
			c.JSON(http.StatusInternalServerError, gin.H{"error": e.Err.Error()})
			return
		case *myerror.BadRequestError:
			c.JSON(http.StatusBadRequest, gin.H{"error": e.Err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.SetCookie("jwt", singedString, 60*60*24, "/", "localhost", false, true)

	c.JSON(http.StatusOK, &response{
		ID:       int64(user.ID),
		Username: user.UserName,
	})
}

// HandlerLogout implements Handler.
func (*handler) HandlerLogout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
