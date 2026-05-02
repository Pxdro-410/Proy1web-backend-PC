package db

import "log"

func Migrate() {
	createPilotosTable := `
	CREATE TABLE IF NOT EXISTS pilotos (
		id             SERIAL PRIMARY KEY,
		name           VARCHAR(150) NOT NULL,
		team           VARCHAR(150) NOT NULL,
		nationality    VARCHAR(100) NOT NULL,
		number         INT NOT NULL,
		championships  INT NOT NULL DEFAULT 0,
		description    TEXT,
		image_path     VARCHAR(300),
		created_at     TIMESTAMP DEFAULT NOW()
	);`

	createRatingsTable := `
	CREATE TABLE IF NOT EXISTS ratings (
		id         SERIAL PRIMARY KEY,
		piloto_id  INT NOT NULL REFERENCES pilotos(id) ON DELETE CASCADE,
		score      INT NOT NULL CHECK (score >= 1 AND score <= 10),
		comment    TEXT,
		created_at TIMESTAMP DEFAULT NOW()
	);`

	if _, err := DB.Exec(createPilotosTable); err != nil {
		log.Fatalf("error creating pilotos table: %v", err)
	}

	if _, err := DB.Exec(createRatingsTable); err != nil {
		log.Fatalf("error creating ratings table: %v", err)
	}

	log.Println("migrations applied")
}
