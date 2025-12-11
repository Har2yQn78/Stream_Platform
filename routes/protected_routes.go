package routes

import (
	controller "github.com/Har2yQn78/Stream_Platform/controllers"
	"github.com/Har2yQn78/Stream_Platform/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine) {
	router.Use(middleware.AuthMiddleware())

	router.GET("/movie/:imdb_id", controller.GetMovieById())
	router.POST("/addmovie", controller.AddMovie())
}
