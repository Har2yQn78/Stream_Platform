package main

import (
	"net/http"

	controller "github.com/Har2yQn78/Stream_Platform/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})
	router.GET("/movies", controller.GetMovies())
	router.GET("/movie/:imdb_id", controller.GetMovieById())

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
