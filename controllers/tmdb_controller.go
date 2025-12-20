package controllers

import (
	"net/http"
	"strconv"

	"github.com/Har2yQn78/Stream_Platform/services"
	"github.com/gin-gonic/gin"
)

var tmdbService = services.NewTMDBService()

func SearchMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("query")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
			return
		}

		page := 1
		if pageStr := c.Query("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}

		results, err := tmdbService.SearchMovies(query, page)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Add full image URLs
		for i := range results.Results {
			results.Results[i].PosterPath = tmdbService.GetFullPosterURL(results.Results[i].PosterPath, "w500")
			results.Results[i].BackdropPath = tmdbService.GetFullBackdropURL(results.Results[i].BackdropPath, "w1280")
		}

		c.JSON(http.StatusOK, results)
	}
}

func SearchTV() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("query")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
			return
		}

		page := 1
		if pageStr := c.Query("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}

		results, err := tmdbService.SearchTV(query, page)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for i := range results.Results {
			results.Results[i].PosterPath = tmdbService.GetFullPosterURL(results.Results[i].PosterPath, "w500")
			results.Results[i].BackdropPath = tmdbService.GetFullBackdropURL(results.Results[i].BackdropPath, "w1280")
		}

		c.JSON(http.StatusOK, results)
	}
}

func GetMovieDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		tmdbIDStr := c.Param("tmdb_id")
		tmdbID, err := strconv.Atoi(tmdbIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid TMDB ID"})
			return
		}

		details, err := tmdbService.GetMovieDetails(tmdbID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		details.PosterPath = tmdbService.GetFullPosterURL(details.PosterPath, "w500")
		details.BackdropPath = tmdbService.GetFullBackdropURL(details.BackdropPath, "w1280")

		c.JSON(http.StatusOK, details)
	}
}

func GetTVDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		tmdbIDStr := c.Param("tmdb_id")
		tmdbID, err := strconv.Atoi(tmdbIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid TMDB ID"})
			return
		}

		details, err := tmdbService.GetTVDetails(tmdbID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		details.PosterPath = tmdbService.GetFullPosterURL(details.PosterPath, "w500")
		details.BackdropPath = tmdbService.GetFullBackdropURL(details.BackdropPath, "w1280")

		c.JSON(http.StatusOK, details)
	}
}
