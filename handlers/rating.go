package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Pxdro-410/Proy1web-backend-PC/db"
	"github.com/Pxdro-410/Proy1web-backend-PC/models"
)

// POST /piloto/:id/rating
func CreateRating(w http.ResponseWriter, r *http.Request) {
	// path: /piloto/{id}/rating
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}

	pilotoID, err := parseInt(parts[1])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid piloto id")
		return
	}

	var input models.CreateRatingInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if input.Score < 1 || input.Score > 10 {
		writeError(w, http.StatusBadRequest, "score must be between 1 and 10")
		return
	}

	var exists bool
	db.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM pilotos WHERE id=$1)`, pilotoID).Scan(&exists)
	if !exists {
		writeError(w, http.StatusNotFound, "piloto not found")
		return
	}

	var rating models.Rating
	err = db.DB.QueryRow(
		`INSERT INTO ratings (piloto_id, score, comment) VALUES ($1,$2,$3)
		 RETURNING id, piloto_id, score, comment, created_at`,
		pilotoID, input.Score, input.Comment,
	).Scan(&rating.ID, &rating.PilotoID, &rating.Score, &rating.Comment, &rating.CreatedAt)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "error saving rating")
		return
	}

	writeJSON(w, http.StatusCreated, rating)
}

// GET /piloto/:id/rating
func GetRatings(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}

	pilotoID, err := parseInt(parts[1])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid piloto id")
		return
	}

	var exists bool
	db.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM pilotos WHERE id=$1)`, pilotoID).Scan(&exists)
	if !exists {
		writeError(w, http.StatusNotFound, "piloto not found")
		return
	}

	rows, err := db.DB.Query(
		`SELECT id, piloto_id, score, comment, created_at FROM ratings WHERE piloto_id=$1 ORDER BY created_at DESC`,
		pilotoID,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error fetching ratings")
		return
	}
	defer rows.Close()

	ratings := []models.Rating{}
	var totalScore int
	for rows.Next() {
		var rat models.Rating
		if err := rows.Scan(&rat.ID, &rat.PilotoID, &rat.Score, &rat.Comment, &rat.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "error reading ratings")
			return
		}
		totalScore += rat.Score
		ratings = append(ratings, rat)
	}

	avg := 0.0
	if len(ratings) > 0 {
		avg = float64(totalScore) / float64(len(ratings))
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"ratings": ratings,
		"average": avg,
		"count":   len(ratings),
	})
}

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
