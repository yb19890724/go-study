package router

import (
	. "controllers"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	router := gin.Default()

	router.GET("/users", Users)

	router.POST("/user", UserStore)

	router.PUT("/user/:id", UserUpdate)

	router.DELETE("/user/:id", UserDestroy)

	router.GET("/pets", Pets)

	return router
}
