package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Pxdro-410/Proy1web-backend-PC/db"
	"github.com/Pxdro-410/Proy1web-backend-PC/models"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, models.ErrorResponse{Error: msg})
}

func idFromPath(r *http.Request, prefix string) (int, error) {
	raw := strings.TrimPrefix(r.URL.Path, prefix)
	raw = strings.Split(raw, "/")[0]
	return strconv.Atoi(raw)
}

// GET /piloto
func GetSeries(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	search := q.Get("q")
	sort := q.Get("sort")
	order := strings.ToUpper(q.Get("order"))

	allowedSorts := map[string]bool{"name": true, "team": true, "number": true, "championships": true, "created_at": true}
	if !allowedSorts[sort] {
		sort = "id"
	}
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	baseQuery := `SELECT id, name, team, nationality, number, championships, description, COALESCE(image_path,'') , created_at FROM pilotos`
	countQuery := `SELECT COUNT(*) FROM pilotos`

	var args []interface{}
	argIdx := 1

	if search != "" {
		where := fmt.Sprintf(" WHERE name ILIKE $%d", argIdx)
		baseQuery += where
		countQuery += where
		args = append(args, "%"+search+"%")
		argIdx++
	}

	var total int
	if err := db.DB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		writeError(w, http.StatusInternalServerError, "error counting records")
		return
	}

	baseQuery += fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", sort, order, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(baseQuery, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error fetching series")
		return
	}
	defer rows.Close()

	pilotos := []models.Piloto{}
	for rows.Next() {
		var p models.Piloto
		if err := rows.Scan(&p.ID, &p.Name, &p.Team, &p.Nationality, &p.Number, &p.Championships, &p.Description, &p.ImagePath, &p.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "error reading records")
			return
		}
		pilotos = append(pilotos, p)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Data:       pilotos,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GET /piloto/:id
func GetSeriesByID(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r, "/piloto/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var p models.Piloto
	err = db.DB.QueryRow(
		`SELECT id, name, team, nationality, number, championships, description, COALESCE(image_path,''), created_at FROM pilotos WHERE id=$1`, id,
	).Scan(&p.ID, &p.Name, &p.Team, &p.Nationality, &p.Number, &p.Championships, &p.Description, &p.ImagePath, &p.CreatedAt)

	if err != nil {
		writeError(w, http.StatusNotFound, "piloto not found")
		return
	}

	writeJSON(w, http.StatusOK, p)
}

// POST /piloto
func CreateSeries(w http.ResponseWriter, r *http.Request) {
	var input models.CreatePilotoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if strings.TrimSpace(input.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if strings.TrimSpace(input.Team) == "" {
		writeError(w, http.StatusBadRequest, "team is required")
		return
	}
	if strings.TrimSpace(input.Nationality) == "" {
		writeError(w, http.StatusBadRequest, "nationality is required")
		return
	}
	if input.Number < 1 || input.Number > 99 {
		writeError(w, http.StatusBadRequest, "number must be between 1 and 99")
		return
	}
	if input.Championships < 0 {
		writeError(w, http.StatusBadRequest, "championships cannot be negative")
		return
	}

	var p models.Piloto
	err := db.DB.QueryRow(
		`INSERT INTO pilotos (name, team, nationality, number, championships, description)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 RETURNING id, name, team, nationality, number, championships, description, COALESCE(image_path,''), created_at`,
		input.Name, input.Team, input.Nationality, input.Number, input.Championships, input.Description,
	).Scan(&p.ID, &p.Name, &p.Team, &p.Nationality, &p.Number, &p.Championships, &p.Description, &p.ImagePath, &p.CreatedAt)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "error creating piloto")
		return
	}

	writeJSON(w, http.StatusCreated, p)
}

// PUT /piloto/:id
func UpdateSeries(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r, "/piloto/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input models.UpdatePilotoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if strings.TrimSpace(input.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if strings.TrimSpace(input.Team) == "" {
		writeError(w, http.StatusBadRequest, "team is required")
		return
	}
	if strings.TrimSpace(input.Nationality) == "" {
		writeError(w, http.StatusBadRequest, "nationality is required")
		return
	}
	if input.Number < 1 || input.Number > 99 {
		writeError(w, http.StatusBadRequest, "number must be between 1 and 99")
		return
	}
	if input.Championships < 0 {
		writeError(w, http.StatusBadRequest, "championships cannot be negative")
		return
	}

	var p models.Piloto
	err = db.DB.QueryRow(
		`UPDATE pilotos SET name=$1, team=$2, nationality=$3, number=$4, championships=$5, description=$6
		 WHERE id=$7
		 RETURNING id, name, team, nationality, number, championships, description, COALESCE(image_path,''), created_at`,
		input.Name, input.Team, input.Nationality, input.Number, input.Championships, input.Description, id,
	).Scan(&p.ID, &p.Name, &p.Team, &p.Nationality, &p.Number, &p.Championships, &p.Description, &p.ImagePath, &p.CreatedAt)

	if err != nil {
		writeError(w, http.StatusNotFound, "piloto not found")
		return
	}

	writeJSON(w, http.StatusOK, p)
}

// DELETE /piloto/:id
func DeleteSeries(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r, "/piloto/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	res, err := db.DB.Exec(`DELETE FROM pilotos WHERE id=$1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error deleting piloto")
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		writeError(w, http.StatusNotFound, "piloto not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// POST /piloto/:id/image
func UploadImage(w http.ResponseWriter, r *http.Request) {
	id, err := idFromPath(r, "/piloto/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	// 1 MB limit
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "image too large, max 1MB")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		writeError(w, http.StatusBadRequest, "field 'image' is required")
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		writeError(w, http.StatusBadRequest, "only jpg, png and webp images are allowed")
		return
	}

	filename := fmt.Sprintf("%d_%d%s", id, time.Now().Unix(), ext)
	dst := filepath.Join("uploads", filename)

	out, err := os.Create(dst)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error saving image")
		return
	}
	defer out.Close()
	io.Copy(out, file)

	imagePath := "/uploads/" + filename
	_, err = db.DB.Exec(`UPDATE pilotos SET image_path=$1 WHERE id=$2`, imagePath, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "piloto not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"image_path": imagePath})
}
