package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Har2yQn78/Stream_Platform/database"
	"github.com/Har2yQn78/Stream_Platform/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")

var validate = validator.New()

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movies []models.Movie

		cursor, err := movieCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, movies)
	}
}

func GetMovieById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")

		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "movie id is empty"})
			return
		}
		var movie models.Movie

		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movie models.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if err := validate.Struct(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		movie.Reviews = []models.Review{}
		movie.Ratings = []models.Rating{}
		movie.AverageRating = 0.0
		movie.TotalRatings = 0

		result, err := movieCollection.InsertOne(ctx, movie)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func AddReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "movie id is empty"})
			return
		}

		userID, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		var reviewRequest models.AddReviewRequest
		if err := c.ShouldBindJSON(&reviewRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := validate.Struct(&reviewRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		userCollection := database.OpenCollection("users")
		err := userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find user"})
			return
		}

		review := models.Review{
			ReviewID:  bson.NewObjectID().Hex(),
			UserID:    userID.(string),
			UserName:  user.FirstName + " " + user.LastName,
			Comment:   reviewRequest.Comment,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		var movie models.Movie
		err = movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		for _, r := range movie.Reviews {
			if r.UserID == userID.(string) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "You have already reviewed this movie. Use update endpoint to modify your review."})
				return
			}
		}

		update := bson.M{
			"$push": bson.M{
				"reviews": review,
			},
		}

		result, err := movieCollection.UpdateOne(ctx, bson.M{"imdb_id": movieID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Review added successfully",
			"review":  review,
			"result":  result,
		})
	}
}

func UpdateReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")
		reviewID := c.Param("review_id")

		if movieID == "" || reviewID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "movie id or review id is empty"})
			return
		}

		userID, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		var updateRequest models.UpdateReviewRequest
		if err := c.ShouldBindJSON(&updateRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := validate.Struct(&updateRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var movie models.Movie
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		found := false
		for _, review := range movie.Reviews {
			if review.ReviewID == reviewID && review.UserID == userID.(string) {
				found = true
				break
			}
		}

		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "Review not found or you don't have permission to update it"})
			return
		}

		// Update the review
		update := bson.M{
			"$set": bson.M{
				"reviews.$[elem].comment":    updateRequest.Comment,
				"reviews.$[elem].updated_at": time.Now(),
			},
		}

		opts := options.UpdateOne().SetArrayFilters([]interface{}{
			bson.M{
				"elem.review_id": reviewID,
				"elem.user_id":   userID.(string),
			},
		})

		result, err := movieCollection.UpdateOne(ctx, bson.M{"imdb_id": movieID}, update, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Review not found or you don't have permission to update it"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Review updated successfully"})
	}
}

func DeleteReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")
		reviewID := c.Param("review_id")

		if movieID == "" || reviewID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "movie id or review id is empty"})
			return
		}

		userID, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		role, _ := c.Get("role")
		var movie models.Movie
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		canDelete := false
		for _, review := range movie.Reviews {
			if review.ReviewID == reviewID {
				if review.UserID == userID.(string) || role == "ADMIN" {
					canDelete = true
					break
				}
			}
		}

		if !canDelete {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this review"})
			return
		}

		update := bson.M{
			"$pull": bson.M{
				"reviews": bson.M{"review_id": reviewID},
			},
		}

		result, err := movieCollection.UpdateOne(ctx, bson.M{"imdb_id": movieID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Review deleted successfully",
			"result":  result,
		})
	}
}

func AddRating() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "movie id is empty"})
			return
		}

		userID, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		var ratingRequest models.AddRatingRequest
		if err := c.ShouldBindJSON(&ratingRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := validate.Struct(&ratingRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var movie models.Movie
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}
		userRatedIndex := -1
		for i, r := range movie.Ratings {
			if r.UserID == userID.(string) {
				userRatedIndex = i
				break
			}
		}

		var update bson.M
		var newAverageRating float64
		var newTotalRatings int

		if userRatedIndex != -1 {
			oldRating := movie.Ratings[userRatedIndex].Rating
			totalSum := movie.AverageRating * float64(movie.TotalRatings)
			totalSum = totalSum - oldRating + ratingRequest.Rating
			newAverageRating = totalSum / float64(movie.TotalRatings)
			newTotalRatings = movie.TotalRatings

			update = bson.M{
				"$set": bson.M{
					"ratings.$[elem].rating":     ratingRequest.Rating,
					"ratings.$[elem].created_at": time.Now(),
					"average_rating":             newAverageRating,
				},
			}

			opts := options.UpdateOne().SetArrayFilters([]interface{}{
				bson.M{"elem.user_id": userID.(string)},
			})

			result, err := movieCollection.UpdateOne(ctx, bson.M{"imdb_id": movieID}, update, opts)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message":        "Rating updated successfully",
				"average_rating": newAverageRating,
				"total_ratings":  newTotalRatings,
				"result":         result,
			})
			return
		}

		rating := models.Rating{
			UserID:    userID.(string),
			Rating:    ratingRequest.Rating,
			CreatedAt: time.Now(),
		}

		totalSum := movie.AverageRating * float64(movie.TotalRatings)
		totalSum += ratingRequest.Rating
		newTotalRatings = movie.TotalRatings + 1
		newAverageRating = totalSum / float64(newTotalRatings)

		update = bson.M{
			"$push": bson.M{
				"ratings": rating,
			},
			"$set": bson.M{
				"average_rating": newAverageRating,
				"total_ratings":  newTotalRatings,
			},
		}

		result, err := movieCollection.UpdateOne(ctx, bson.M{"imdb_id": movieID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":        "Rating added successfully",
			"average_rating": newAverageRating,
			"total_ratings":  newTotalRatings,
			"result":         result,
		})
	}
}

func GetMovieReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "movie id is empty"})
			return
		}

		var movie models.Movie
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"reviews":        movie.Reviews,
			"average_rating": movie.AverageRating,
			"total_ratings":  movie.TotalRatings,
		})
	}
}
