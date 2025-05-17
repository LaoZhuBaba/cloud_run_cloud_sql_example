package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

// visitData is used to pass data to the HTML template.
type visitData struct {
	RecentVisits []visit
}

// visit contains a single row from the visits table in the database.
// Each visit includes a timestamp.
type visit struct {
	VisitTime time.Time
}

func connect(ctx context.Context) (*pgx.Conn, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Fatal Error in connect_connector.go: %s environment variable not set.\n", k)
		}
		return v
	}
	var (
		dbUser   = mustGetenv("DB_USER")
		dbPwd    = mustGetenv("DB_PASS")
		dbName   = mustGetenv("DB_NAME")
		dbIPAddr = mustGetenv("DB_IP_ADDRESS")
		dbIPPort = mustGetenv("DB_IP_PORT")
	)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPwd, dbIPAddr, dbIPPort, dbName)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	createVisits := `CREATE TABLE IF NOT EXISTS visits (
		id SERIAL NOT NULL,
		created_at timestamp NOT NULL,
		PRIMARY KEY (id)
	  );`
	_, err = conn.Exec(ctx, createVisits)
	if err != nil {
		log.Fatalf("unable to create table: %s", err)
	}
	return conn, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)

	ctx := context.Background()

	db, err := connect(ctx)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		// Insert current visit
		_, err := db.Exec(ctx, "INSERT INTO visits(created_at) VALUES(NOW())")
		if err != nil {
			log.Fatalf("unable to save visit: %v", err)
		}

		// Get the last 5 visits
		rows, err := db.Query(ctx, "SELECT created_at FROM visits ORDER BY created_at DESC LIMIT 5")
		if err != nil {
			log.Fatalf("DB.Query: %v", err)
		}
		defer rows.Close()

		var visits []visit
		for rows.Next() {
			var visitTime time.Time
			err := rows.Scan(&visitTime)
			if err != nil {
				log.Fatalf("Rows.Scan: %v", err)
			}
			visits = append(visits, visit{VisitTime: visitTime})
		}
		response, err := json.Marshal(visitData{RecentVisits: visits})
		if err != nil {
			log.Fatalf("renderIndex: failed to parse totals with json.Marshal: %v", err)
		}
		w.Write(response)
	})
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
