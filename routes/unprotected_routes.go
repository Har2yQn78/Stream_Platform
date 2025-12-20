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
	// Auth routes
	router.POST("/register", controller.RegisterUser())
	router.POST("/login", controller.LoginUser(database.Client))

	router.GET("/media", controller.GetAllMedia())
	router.GET("/media/:tmdb_id", controller.GetMediaByTMDBID())
	router.GET("/media/:tmdb_id/reviews", controller.GetMediaReviews())
	router.GET("/media/:tmdb_id/comments", controller.GetMediaComments())
}
