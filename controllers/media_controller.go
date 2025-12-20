package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Har2yQn78/Stream_Platform/database"
	"github.com/Har2yQn78/Stream_Platform/models"
	"github.com/Har2yQn78/Stream_Platform/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var mediaCollection *mongo.Collection = database.OpenCollection("media")
var mediaValidator = validator.New()
var mediaTmdbService = services.NewTMDBService()

func GetAllMedia() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var media []models.Media

		filter := bson.M{}
		if mediaType := c.Query("type"); mediaType != "" {
			filter["media_type"] = mediaType
		}

		cursor, err := mediaCollection.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &media); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, media)
	}
}

func GetMediaByTMDBID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		tmdbIDStr := c.Param("tmdb_id")
		tmdbID, err := strconv.Atoi(tmdbIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid TMDB ID"})
			return
		}

		var media models.Media
		err = mediaCollection.FindOne(ctx, bson.M{"tmdb_id": tmdbID}).Decode(&media)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
			return
		}

		c.JSON(http.StatusOK, media)
	}
}

func AddMedia() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userID, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		var request models.AddMediaRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := mediaValidator.Struct(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		count, err := mediaCollection.CountDocuments(ctx, bson.M{"tmdb_id": request.TMDBID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Media already exists in database"})
			return
		}

		var media models.Media

		if request.MediaType == models.MediaTypeMovie {
			details, err := mediaTmdbService.GetMovieDetails(request.TMDBID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Could not fetch movie details from TMDB: " + err.Error()})
				return
			}

			genres := make([]models.Genre, len(details.Genres))
			for i, g := range details.Genres {
				genres[i] = models.Genre{
					GenreID:   g.ID,
					GenreName: g.Name,
				}
			}

			media = models.Media{
				TMDBID:       details.ID,
				ImdbID:       details.ImdbID,
				MediaType:    models.MediaTypeMovie,
				Title:        details.Title,
				Overview:     details.Overview,
				PosterPath:   mediaTmdbService.GetFullPosterURL(details.PosterPath, "w500"),
				BackdropPath: mediaTmdbService.GetFullBackdropURL(details.BackdropPath, "w1280"),
				VideoURL:     request.VideoURL,
				ReleaseDate:  details.ReleaseDate,
				Genres:       genres,
				Runtime:      details.Runtime,
			}
		} else if request.MediaType == models.MediaTypeTV {
			details, err := mediaTmdbService.GetTVDetails(request.TMDBID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Could not fetch TV details from TMDB: " + err.Error()})
				return
			}

			genres := make([]models.Genre, len(details.Genres))
			for i, g := range details.Genres {
				genres[i] = models.Genre{
					GenreID:   g.ID,
					GenreName: g.Name,
				}
			}

			media = models.Media{
				TMDBID:           details.ID,
				MediaType:        models.MediaTypeTV,
				Title:            details.Name,
				Overview:         details.Overview,
				PosterPath:       mediaTmdbService.GetFullPosterURL(details.PosterPath, "w500"),
				BackdropPath:     mediaTmdbService.GetFullBackdropURL(details.BackdropPath, "w1280"),
				VideoURL:         request.VideoURL,
				ReleaseDate:      details.FirstAirDate,
				Genres:           genres,
				NumberOfSeasons:  details.NumberOfSeasons,
				NumberOfEpisodes: details.NumberOfEpisodes,
				InProduction:     details.InProduction,
			}
		}

		media.Reviews = []models.Review{}
		media.Comments = []models.Comment{}
		media.Ratings = []models.Rating{}
		media.AverageRating = 0.0
		media.TotalRatings = 0
		media.AddedBy = userID.(string)
		media.CreatedAt = time.Now()
		media.UpdatedAt = time.Now()

		result, err := mediaCollection.InsertOne(ctx, media)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":   "Media added successfully",
			"insert_id": result.InsertedID,
			"media":     media,
		})
	}
}

func AddMediaReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		tmdbIDStr := c.Param("tmdb_id")
		tmdbID, err := strconv.Atoi(tmdbIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid TMDB ID"})
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

		if err := mediaValidator.Struct(&reviewRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		userCollection := database.OpenCollection("users")
		err = userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find user"})
			return
		}

		var media models.Media
		err = mediaCollection.FindOne(ctx, bson.M{"tmdb_id": tmdbID}).Decode(&media)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
			return
		}

		for _, r := range media.Reviews {
			if r.UserID == userID.(string) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "You have already reviewed this media. Use update endpoint to modify your review."})
				return
			}
		}

		review := models.Review{
			ReviewID:  bson.NewObjectID().Hex(),
			UserID:    userID.(string),
			UserName:  user.FirstName + " " + user.LastName,
			Comment:   reviewRequest.Comment,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		update := bson.M{
			"$push": bson.M{"reviews": review},
			"$set":  bson.M{"updated_at": time.Now()},
		}

		result, err := mediaCollection.UpdateOne(ctx, bson.M{"tmdb_id": tmdbID}, update)
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

func AddMediaComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		tmdbIDStr := c.Param("tmdb_id")
		tmdbID, err := strconv.Atoi(tmdbIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid TMDB ID"})
			return
		}

		userID, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		var commentRequest models.AddCommentRequest
		if err := c.ShouldBindJSON(&commentRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := mediaValidator.Struct(&commentRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		userCollection := database.OpenCollection("users")
		err = userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find user"})
			return
		}

		var media models.Media
		err = mediaCollection.FindOne(ctx, bson.M{"tmdb_id": tmdbID}).Decode(&media)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
			return
		}

		comment := models.Comment{
			CommentID: bson.NewObjectID().Hex(),
			UserID:    userID.(string),
			UserName:  user.FirstName + " " + user.LastName,
			Content:   commentRequest.Content,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		update := bson.M{
			"$push": bson.M{"comments": comment},
			"$set":  bson.M{"updated_at": time.Now()},
		}

		result, err := mediaCollection.UpdateOne(ctx, bson.M{"tmdb_id": tmdbID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Comment added successfully",
			"comment": comment,
			"result":  result,
		})
	}
}

func DeleteMediaComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		tmdbIDStr := c.Param("tmdb_id")
		tmdbID, err := strconv.Atoi(tmdbIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid TMDB ID"})
			return
		}

		commentID := c.Param("comment_id")
		if commentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "comment_id is required"})
			return
		}

		userID, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		role, _ := c.Get("role")

		var media models.Media
		err = mediaCollection.FindOne(ctx, bson.M{"tmdb_id": tmdbID}).Decode(&media)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
			return
		}

		canDelete := false
		for _, comment := range media.Comments {
			if comment.CommentID == commentID {
				if comment.UserID == userID.(string) || role == "ADMIN" {
					canDelete = true
					break
				}
			}
		}

		if !canDelete {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this comment"})
			return
		}

		update := bson.M{
			"$pull": bson.M{"comments": bson.M{"comment_id": commentID}},
			"$set":  bson.M{"updated_at": time.Now()},
		}

		result, err := mediaCollection.UpdateOne(ctx, bson.M{"tmdb_id": tmdbID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Comment deleted successfully",
			"result":  result,
		})
	}
}

func GetMediaComments() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		tmdbIDStr := c.Param("tmdb_id")
		tmdbID, err := strconv.Atoi(tmdbIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid TMDB ID"})
			return
		}

		var media models.Media
		err = mediaCollection.FindOne(ctx, bson.M{"tmdb_id": tmdbID}).Decode(&media)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"comments": media.Comments,
			"total":    len(media.Comments),
		})
	}
}

func AddMediaRating() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		tmdbIDStr := c.Param("tmdb_id")
		tmdbID, err := strconv.Atoi(tmdbIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid TMDB ID"})
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

		if err := mediaValidator.Struct(&ratingRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var media models.Media
		err = mediaCollection.FindOne(ctx, bson.M{"tmdb_id": tmdbID}).Decode(&media)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
			return
		}

		userRatedIndex := -1
		for i, r := range media.Ratings {
			if r.UserID == userID.(string) {
				userRatedIndex = i
				break
			}
		}

		var update bson.M
		var newAverageRating float64
		var newTotalRatings int

		if userRatedIndex != -1 {
			oldRating := media.Ratings[userRatedIndex].Rating
			totalSum := media.AverageRating * float64(media.TotalRatings)
			totalSum = totalSum - oldRating + ratingRequest.Rating
			newAverageRating = totalSum / float64(media.TotalRatings)
			newTotalRatings = media.TotalRatings

			update = bson.M{
				"$set": bson.M{
					"ratings.$[elem].rating":     ratingRequest.Rating,
					"ratings.$[elem].created_at": time.Now(),
					"average_rating":             newAverageRating,
					"updated_at":                 time.Now(),
				},
			}

			opts := options.UpdateOne().SetArrayFilters([]interface{}{
				bson.M{"elem.user_id": userID.(string)},
			})

			result, err := mediaCollection.UpdateOne(ctx, bson.M{"tmdb_id": tmdbID}, update, opts)
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

		totalSum := media.AverageRating * float64(media.TotalRatings)
		totalSum += ratingRequest.Rating
		newTotalRatings = media.TotalRatings + 1
		newAverageRating = totalSum / float64(newTotalRatings)

		update = bson.M{
			"$push": bson.M{"ratings": rating},
			"$set": bson.M{
				"average_rating": newAverageRating,
				"total_ratings":  newTotalRatings,
				"updated_at":     time.Now(),
			},
		}

		result, err := mediaCollection.UpdateOne(ctx, bson.M{"tmdb_id": tmdbID}, update)
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

func GetMediaReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		tmdbIDStr := c.Param("tmdb_id")
		tmdbID, err := strconv.Atoi(tmdbIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid TMDB ID"})
			return
		}

		var media models.Media
		err = mediaCollection.FindOne(ctx, bson.M{"tmdb_id": tmdbID}).Decode(&media)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"reviews":        media.Reviews,
			"average_rating": media.AverageRating,
			"total_ratings":  media.TotalRatings,
		})
	}
}
