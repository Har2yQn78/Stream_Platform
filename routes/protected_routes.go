package routes

import (
	controller "github.com/Har2yQn78/Stream_Platform/controllers"
	"github.com/Har2yQn78/Stream_Platform/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine) {
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/addmovie", controller.AddMovie())
		protected.POST("/movie/:imdb_id/review", controller.AddReview())
		protected.PUT("/movie/:imdb_id/review/:review_id", controller.UpdateReview())
		protected.DELETE("/movie/:imdb_id/review/:review_id", controller.DeleteReview())
		protected.POST("/movie/:imdb_id/rating", controller.AddRating())

		protected.POST("/media", controller.AddMedia())
		protected.POST("/media/:tmdb_id/review", controller.AddMediaReview())
		protected.POST("/media/:tmdb_id/comment", controller.AddMediaComment())
		protected.DELETE("/media/:tmdb_id/comment/:comment_id", controller.DeleteMediaComment())
		protected.POST("/media/:tmdb_id/rating", controller.AddMediaRating())
	}
}
