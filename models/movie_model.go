package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Genre struct {
	GenreID   int    `bson:"genre_id" json:"genre_id" validate:"required"`
	GenreName string `bson:"genre_name" json:"genre_name" validate:"required,min=2,max=100"`
}

type Ranking struct {
	RankingValue int    `bson:"ranking_value" json:"ranking_value" validate:"required"`
	RankingName  string `bson:"ranking_name" json:"ranking_name" validate:"required"`
}

type Review struct {
	ReviewID  string    `bson:"review_id" json:"review_id"`
	UserID    string    `bson:"user_id" json:"user_id" validate:"required"`
	UserName  string    `bson:"user_name" json:"user_name" validate:"required"`
	Comment   string    `bson:"comment" json:"comment" validate:"required,min=10,max=1000"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type Rating struct {
	UserID    string    `bson:"user_id" json:"user_id" validate:"required"`
	Rating    float64   `bson:"rating" json:"rating" validate:"required,min=0,max=10"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type Movie struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ImdbID        string        `bson:"imdb_id" json:"imdb_id" validate:"required"`
	Title         string        `bson:"title" json:"title" validate:"required,min=2,max=500"`
	PosterPath    string        `bson:"poster_path" json:"poster_path" validate:"required,url"`
	YouTubeID     string        `bson:"youtube_id" json:"youtube_id" validate:"required"`
	Genre         []Genre       `bson:"genre" json:"genre" validate:"required,dive"`
	Reviews       []Review      `bson:"reviews" json:"reviews"`
	Ratings       []Rating      `bson:"ratings" json:"ratings"`
	AverageRating float64       `bson:"average_rating" json:"average_rating"`
	TotalRatings  int           `bson:"total_ratings" json:"total_ratings"`
	Ranking       Ranking       `bson:"ranking" json:"ranking" validate:"required"`
}

type AddReviewRequest struct {
	Comment string `json:"comment" validate:"required,min=10,max=1000"`
}

type AddRatingRequest struct {
	Rating float64 `json:"rating" validate:"required,min=0,max=10"`
}

type UpdateReviewRequest struct {
	Comment string `json:"comment" validate:"required,min=10,max=1000"`
}
