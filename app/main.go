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

// func connectWithConnector() (*sql.DB, error) {
// 	mustGetenv := func(k string) string {
// 		v := os.Getenv(k)
// 		if v == "" {
// 			log.Fatalf("Fatal Error in connect_connector.go: %s environment variable not set.\n", k)
// 		}
// 		return v
// 	}
// 	// Note: Saving credentials in environment variables is convenient, but not
// 	// secure - consider a more secure solution such as
// 	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
// 	// keep passwords and other secrets safe.
// 	var (
// 		dbUser                 = mustGetenv("DB_USER")                  // e.g. 'my-db-user'
// 		dbPwd                  = mustGetenv("DB_PASS")                  // e.g. 'my-db-password'
// 		dbName                 = mustGetenv("DB_NAME")                  // e.g. 'my-database'
// 		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
// 	)

// 	dsn := fmt.Sprintf("user=%s password=%s database=%s", dbUser, dbPwd, dbName)
// 	config, err := pgx.ParseConfig(dsn)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var opts []cloudsqlconn.Option

// 	// WithLazyRefresh() Option is used to perform refresh
// 	// when needed, rather than on a scheduled interval.
// 	// This is recommended for serverless environments to
// 	// avoid background refreshes from throttling CPU.
// 	opts = append(opts, cloudsqlconn.WithLazyRefresh())
// 	d, err := cloudsqlconn.NewDialer(context.Background(), opts...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Use the Cloud SQL connector to handle connecting to the instance.
// 	// This approach does *NOT* require the Cloud SQL proxy.
// 	config.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
// 		conn, err := d.Dial(ctx, instanceConnectionName, cloudsqlconn.WithPSC())
// 		if err != nil {
// 			log.Fatalf("Dial failed with error: %v", err)
// 		}
// 		return conn, nil
// 	}
// 	dbURI := stdlib.RegisterConnConfig(config)
// 	dbPool, err := sql.Open("pgx", dbURI)
// 	if err != nil {
// 		return nil, fmt.Errorf("sql.Open: %w", err)
// 	}
// 	return dbPool, nil
// }

func connect(ctx context.Context) (*pgx.Conn, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Fatal Error in connect_connector.go: %s environment variable not set.\n", k)
		}
		return v
	}
	// Note: Saving credentials in environment variables is convenient, but not
	// secure - consider a more secure solution such as
	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
	// keep passwords and other secrets safe.
	var (
		dbUser   = mustGetenv("DB_USER") // e.g. 'my-db-user'
		dbPwd    = mustGetenv("DB_PASS") // e.g. 'my-db-password'
		dbName   = mustGetenv("DB_NAME") // e.g. 'my-database'
		dbIPAddr = mustGetenv("DB_IP_ADDRESS")
		dbIPPort = mustGetenv("DB_IP_PORT")
	)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPwd, dbIPAddr, dbIPPort, dbName)

	log.Printf("connection string: %s", fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, "PW_NOT_SHOWN", dbIPAddr, dbIPPort, dbName))

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
