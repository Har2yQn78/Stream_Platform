package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
