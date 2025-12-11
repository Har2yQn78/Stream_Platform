package routes

import (
	controller "github.com/Har2yQn78/Stream_Platform/controllers"
	"github.com/Har2yQn78/Stream_Platform/database"
	"github.com/gin-gonic/gin"
)

func SetupUnProtectedRoutes(router *gin.Engine) {

	router.GET("/movies", controller.GetMovies())
	router.POST("/register", controller.RegisterUser())
	router.POST("/login", controller.LoginUser(database.Client))
}
