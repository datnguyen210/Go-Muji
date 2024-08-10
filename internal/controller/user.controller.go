package controller

import (
	"net/http"

	"github.com/datnguyen210/go-muji/internal/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

func (uc *UserController) GetUserByID(c *gin.Context) {
	// id := c.Query("id")
	c.JSON(http.StatusOK, gin.H{
		"name": uc.userService.ReadUserInfo(),
	})
}
