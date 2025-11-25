package main

import (
	controller "github.com/Har2yQn78/Stream_Platform/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})
	router.GET("/movies", controller.GetMovies())

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
