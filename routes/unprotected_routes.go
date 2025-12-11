package routes

import (
	controller "github.com/Har2yQn78/Stream_Platform/controllers"
	"github.com/Har2yQn78/Stream_Platform/database"
	"github.com/gin-gonic/gin"
)

func SetupUnProtectedRoutes(router *gin.Engine) {
	router.GET("/movies", controller.GetMovies())
	router.GET("/movie/:imdb_id", controller.GetMovieById())
	router.GET("/movie/:imdb_id/reviews", controller.GetMovieReviews())
	// Auth points
	router.POST("/register", controller.RegisterUser())
	router.POST("/login", controller.LoginUser(database.Client))
}
