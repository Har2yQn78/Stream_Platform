package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MediaType string

const (
	MediaTypeMovie MediaType = "movie"
	MediaTypeTV    MediaType = "tv"
)

type Comment struct {
	CommentID string    `bson:"comment_id" json:"comment_id"`
	UserID    string    `bson:"user_id" json:"user_id" validate:"required"`
	UserName  string    `bson:"user_name" json:"user_name" validate:"required"`
	Content   string    `bson:"content" json:"content" validate:"required,min=1,max=500"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type Media struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	TMDBID    int           `bson:"tmdb_id" json:"tmdb_id" validate:"required"`
	ImdbID    string        `bson:"imdb_id,omitempty" json:"imdb_id,omitempty"`
	MediaType MediaType     `bson:"media_type" json:"media_type" validate:"required,oneof=movie tv"`

	Title        string `bson:"title" json:"title" validate:"required,min=1,max=500"`
	Overview     string `bson:"overview" json:"overview"`
	PosterPath   string `bson:"poster_path" json:"poster_path"`
	BackdropPath string `bson:"backdrop_path" json:"backdrop_path"`
	VideoURL     string `bson:"video_url" json:"video_url" validate:"required,url"`
	ReleaseDate  string `bson:"release_date" json:"release_date"`

	Genres []Genre `bson:"genres" json:"genres"`

	NumberOfSeasons  int  `bson:"number_of_seasons,omitempty" json:"number_of_seasons,omitempty"`
	NumberOfEpisodes int  `bson:"number_of_episodes,omitempty" json:"number_of_episodes,omitempty"`
	InProduction     bool `bson:"in_production,omitempty" json:"in_production,omitempty"`

	Runtime int `bson:"runtime,omitempty" json:"runtime,omitempty"`

	Reviews       []Review  `bson:"reviews" json:"reviews"`
	Comments      []Comment `bson:"comments" json:"comments"`
	Ratings       []Rating  `bson:"ratings" json:"ratings"`
	AverageRating float64   `bson:"average_rating" json:"average_rating"`
	TotalRatings  int       `bson:"total_ratings" json:"total_ratings"`

	Ranking   Ranking   `bson:"ranking,omitempty" json:"ranking,omitempty"`
	AddedBy   string    `bson:"added_by" json:"added_by"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type AddMediaRequest struct {
	TMDBID    int       `json:"tmdb_id" validate:"required"`
	MediaType MediaType `json:"media_type" validate:"required,oneof=movie tv"`
	VideoURL  string    `json:"video_url" validate:"required,url"`
}

type AddCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=500"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=500"`
}
