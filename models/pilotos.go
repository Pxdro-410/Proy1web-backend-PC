package models

import "time"

type Piloto struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Team        string    `json:"team"`
	Nationality string    `json:"nationality"`
	Number      int       `json:"number"`
	Championships int     `json:"championships"`
	Description string    `json:"description"`
	ImagePath   string    `json:"image_path"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreatePilotoInput struct {
	Name          string `json:"name"`
	Team          string `json:"team"`
	Nationality   string `json:"nationality"`
	Number        int    `json:"number"`
	Championships int    `json:"championships"`
	Description   string `json:"description"`
}

type UpdatePilotoInput struct {
	Name          string `json:"name"`
	Team          string `json:"team"`
	Nationality   string `json:"nationality"`
	Number        int    `json:"number"`
	Championships int    `json:"championships"`
	Description   string `json:"description"`
}

type Rating struct {
	ID        int       `json:"id"`
	PilotoID  int       `json:"piloto_id"`
	Score     int       `json:"score"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateRatingInput struct {
	Score   int    `json:"score"`
	Comment string `json:"comment"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int         `json:"total"`
	TotalPages int         `json:"total_pages"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
