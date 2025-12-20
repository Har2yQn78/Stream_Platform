package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// TMDB Configuration
type TMDBConfig struct {
	APIKey       string
	BaseURL      string
	ImageBaseURL string
}

func GetTMDBConfig() *TMDBConfig {
	return &TMDBConfig{
		APIKey:       os.Getenv("TMDB_API_KEY"),
		BaseURL:      "https://api.themoviedb.org/3",
		ImageBaseURL: "https://image.tmdb.org/t/p",
	}
}

// Genre represents a genre from TMDB
type TMDBGenre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TMDBMovie represents a movie from TMDB search results
type TMDBMovie struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	OriginalTitle string  `json:"original_title"`
	Overview      string  `json:"overview"`
	PosterPath    string  `json:"poster_path"`
	BackdropPath  string  `json:"backdrop_path"`
	ReleaseDate   string  `json:"release_date"`
	VoteAverage   float64 `json:"vote_average"`
	VoteCount     int     `json:"vote_count"`
	Popularity    float64 `json:"popularity"`
	GenreIDs      []int   `json:"genre_ids"`
	Adult         bool    `json:"adult"`
}

// TMDBMovieDetails represents detailed movie info from TMDB
type TMDBMovieDetails struct {
	ID            int         `json:"id"`
	ImdbID        string      `json:"imdb_id"`
	Title         string      `json:"title"`
	OriginalTitle string      `json:"original_title"`
	Overview      string      `json:"overview"`
	PosterPath    string      `json:"poster_path"`
	BackdropPath  string      `json:"backdrop_path"`
	ReleaseDate   string      `json:"release_date"`
	Runtime       int         `json:"runtime"`
	VoteAverage   float64     `json:"vote_average"`
	VoteCount     int         `json:"vote_count"`
	Popularity    float64     `json:"popularity"`
	Genres        []TMDBGenre `json:"genres"`
	Status        string      `json:"status"`
	Tagline       string      `json:"tagline"`
	Adult         bool        `json:"adult"`
}

// TMDBTVShow represents a TV show from TMDB search results
type TMDBTVShow struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	OriginalName  string   `json:"original_name"`
	Overview      string   `json:"overview"`
	PosterPath    string   `json:"poster_path"`
	BackdropPath  string   `json:"backdrop_path"`
	FirstAirDate  string   `json:"first_air_date"`
	VoteAverage   float64  `json:"vote_average"`
	VoteCount     int      `json:"vote_count"`
	Popularity    float64  `json:"popularity"`
	GenreIDs      []int    `json:"genre_ids"`
	OriginCountry []string `json:"origin_country"`
}

// TMDBTVDetails represents detailed TV show info from TMDB
type TMDBTVDetails struct {
	ID               int         `json:"id"`
	Name             string      `json:"name"`
	OriginalName     string      `json:"original_name"`
	Overview         string      `json:"overview"`
	PosterPath       string      `json:"poster_path"`
	BackdropPath     string      `json:"backdrop_path"`
	FirstAirDate     string      `json:"first_air_date"`
	LastAirDate      string      `json:"last_air_date"`
	NumberOfSeasons  int         `json:"number_of_seasons"`
	NumberOfEpisodes int         `json:"number_of_episodes"`
	VoteAverage      float64     `json:"vote_average"`
	VoteCount        int         `json:"vote_count"`
	Popularity       float64     `json:"popularity"`
	Genres           []TMDBGenre `json:"genres"`
	Status           string      `json:"status"`
	Tagline          string      `json:"tagline"`
	Type             string      `json:"type"`
	InProduction     bool        `json:"in_production"`
}

// TMDBSearchResponse represents a paginated search response
type TMDBSearchResponse[T any] struct {
	Page         int `json:"page"`
	Results      []T `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

// TMDBService provides methods to interact with TMDB API
type TMDBService struct {
	config *TMDBConfig
	client *http.Client
}

// NewTMDBService creates a new TMDB service instance
func NewTMDBService() *TMDBService {
	return &TMDBService{
		config: GetTMDBConfig(),
		client: &http.Client{},
	}
}

// makeRequest performs an HTTP request to TMDB API
func (s *TMDBService) makeRequest(endpoint string, params map[string]string) ([]byte, error) {
	reqURL, err := url.Parse(s.config.BaseURL + endpoint)
	if err != nil {
		return nil, err
	}

	query := reqURL.Query()
	query.Set("api_key", s.config.APIKey)
	for key, value := range params {
		query.Set(key, value)
	}
	reqURL.RawQuery = query.Encode()

	resp, err := s.client.Get(reqURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("TMDB API error: %d - %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// SearchMovies searches for movies in TMDB
func (s *TMDBService) SearchMovies(query string, page int) (*TMDBSearchResponse[TMDBMovie], error) {
	if page < 1 {
		page = 1
	}

	body, err := s.makeRequest("/search/movie", map[string]string{
		"query": query,
		"page":  fmt.Sprintf("%d", page),
	})
	if err != nil {
		return nil, err
	}

	var result TMDBSearchResponse[TMDBMovie]
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SearchTV searches for TV shows in TMDB
func (s *TMDBService) SearchTV(query string, page int) (*TMDBSearchResponse[TMDBTVShow], error) {
	if page < 1 {
		page = 1
	}

	body, err := s.makeRequest("/search/tv", map[string]string{
		"query": query,
		"page":  fmt.Sprintf("%d", page),
	})
	if err != nil {
		return nil, err
	}

	var result TMDBSearchResponse[TMDBTVShow]
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetMovieDetails gets detailed information about a movie
func (s *TMDBService) GetMovieDetails(tmdbID int) (*TMDBMovieDetails, error) {
	body, err := s.makeRequest(fmt.Sprintf("/movie/%d", tmdbID), nil)
	if err != nil {
		return nil, err
	}

	var result TMDBMovieDetails
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetTVDetails gets detailed information about a TV show
func (s *TMDBService) GetTVDetails(tmdbID int) (*TMDBTVDetails, error) {
	body, err := s.makeRequest(fmt.Sprintf("/tv/%d", tmdbID), nil)
	if err != nil {
		return nil, err
	}

	var result TMDBTVDetails
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetFullPosterURL constructs full poster URL from path
func (s *TMDBService) GetFullPosterURL(posterPath string, size string) string {
	if posterPath == "" {
		return ""
	}
	if size == "" {
		size = "w500"
	}
	return fmt.Sprintf("%s/%s%s", s.config.ImageBaseURL, size, posterPath)
}

// GetFullBackdropURL constructs full backdrop URL from path
func (s *TMDBService) GetFullBackdropURL(backdropPath string, size string) string {
	if backdropPath == "" {
		return ""
	}
	if size == "" {
		size = "w1280"
	}
	return fmt.Sprintf("%s/%s%s", s.config.ImageBaseURL, size, backdropPath)
}
