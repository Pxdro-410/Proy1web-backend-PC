package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Pxdro-410/Proy1web-backend-PC/db"
	"github.com/Pxdro-410/Proy1web-backend-PC/handlers"
	"github.com/Pxdro-410/Proy1web-backend-PC/middleware"
)

func main() {
	db.Connect()
	db.Migrate()

	mux := http.NewServeMux()

	// GET /piloto  |  POST /piloto
	mux.HandleFunc("/piloto", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetPilotos(w, r)
		case http.MethodPost:
			handlers.CreatePiloto(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// /piloto/:id  |  /piloto/:id/rating
	mux.HandleFunc("/piloto/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if strings.HasSuffix(path, "/rating") {
			switch r.Method {
			case http.MethodGet:
				handlers.GetRatings(w, r)
			case http.MethodPost:
				handlers.CreateRating(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		switch r.Method {
		case http.MethodGet:
			handlers.GetPilotoByID(w, r)
		case http.MethodPut:
			handlers.UpdatePiloto(w, r)
		case http.MethodDelete:
			handlers.DeletePiloto(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server running on :%s", port)
	if err := http.ListenAndServe(":"+port, middleware.CORS(mux)); err != nil {
		log.Fatal(err)
	}
}
