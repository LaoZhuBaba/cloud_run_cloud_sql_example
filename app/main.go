package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/oauth2/google"
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

func getToken(ctx context.Context) string {
	// This scope is a broad generic scope.  It should be possible to restrict this so that
	// the Oauth2 token only supports required SQL access and nothing else.  To do.
	scopes := []string{"https://www.googleapis.com/auth/cloud-platform"}
	creds, err := google.FindDefaultCredentials(ctx, scopes...)
	if err != nil {
		log.Fatal(err)
	}
	token, err := creds.TokenSource.Token()
	if err != nil {
		log.Fatal(err)
	}
	return token.AccessToken
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
		dbName   = mustGetenv("DB_NAME")
		dbIPAddr = mustGetenv("DB_IP_ADDRESS")
		dbIPPort = mustGetenv("DB_IP_PORT")
		dbPwd    = os.Getenv("DB_PASS") // password is not manditory
	)

	// This is hacky way of making the demo code support both a traditional login/password (for classic Postgresql DB accounts),
	// and support a Google IAM service account which can log in using an Oauth2 token as a password.  For Postgresql you have to
	// remove the .gserviceacount.com part of the suffix in order to make this work.  Odd because this isn't required for MS SQL
	// or MySQL.
	if dbPwd == "" {
		dbPwd = getToken(ctx)
		if dbUserShort, ok := strings.CutSuffix(dbUser, ".gserviceaccount.com"); !ok {
			log.Fatal("failed to truncate dbUser")
		} else {
			log.Printf("dbUser shortened okay: %s -> %s", dbUser, dbUserShort)
			dbUser = dbUserShort
		}
	}

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
