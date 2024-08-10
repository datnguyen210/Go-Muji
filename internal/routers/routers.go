package routers

import (
	"net/http"

	"github.com/datnguyen210/go-muji/internal/controller"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		v1.GET("/user/1", controller.NewUserController().GetUserByID)
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "peng",
			})
		})
	}
	return r
}
