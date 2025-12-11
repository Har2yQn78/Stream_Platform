package main

import (
	"net/http"

	"github.com/Har2yQn78/Stream_Platform/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	routes.SetupUnProtectedRoutes(router)

	routes.SetupProtectedRoutes(router)

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
