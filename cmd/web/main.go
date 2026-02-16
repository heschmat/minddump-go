package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	// MySQL driver for Go's database/sql package.
	// The underscore means we import the package solely for its side-effects
	// (i.e., to register the MySQL driver with database/sql (via its init() function)),
	// and we don't directly reference any of its exported identifiers in our code.
	_ "github.com/go-sql-driver/mysql"
	"github.com/heschmat/minddump-go/internal/models"
)

type config struct {
	addr      string
	staticDir string
	dsn       string
}

// a struct to hold the application-wide dependencies
type application struct {
	logger   *slog.Logger
	error    *errorHandler
	snippets *models.SnippetModel
}

func main() {
	// logger -----------------------------------
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // for development
		// AddSource: true,
	}))

	// flags & env variables --------------------
	var cfg config

	dsn := os.Getenv("SNIPPETBOX_DSN")
	if dsn == "" {
		logger.Error("SNIPPETBOX_DSN must be set")
		os.Exit(1)
	}

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.dsn, "dsn", dsn, "MySQL data source name")
	flag.Parse()

	// database connection ----------------------
	logger.Info("connecting to database...")
	db, err := openDB(cfg.dsn)
	if err != nil {
		// logger.Error(err.Error())
		logger.Error("unable to connect to database", "error", err)
		os.Exit(1)
	}

	// make sure to close the database connection pool before the main() function exits
	// N.B. this call **here** is actually superfluous
	// as the app will only be terminated by a signal interrupt (e.g., Ctrl+C or by os.Exit(1)),
	// in both cases the app exists immediately & defferred functions won't be executed.
	// but it's still a good practice to include this call in case the app is terminated in a different way (e.g., by an error in the code or by a panic).
	defer db.Close()

	// application dependencies -----------------
	app := &application{
		logger:   logger,
		error:    &errorHandler{logger: logger},
		snippets: &models.SnippetModel{DB: db},
	}

	// server -----------------------------------
	logger.Info("starting server...", "addr", cfg.addr)

	// params: the TCP network address, the servemux
	err = http.ListenAndServe(cfg.addr, app.routes())
	// N.B. Any error returned by http.ListenAndServe() is always non-nil.
	// log.Fatal(err)
	logger.Error(err.Error())
	os.Exit(1)
}

// a helper function which returns a sql.DB connection pool for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	// sql.Open() does not establish any connections to the database by default.
	// Instead, it just validates the DSN and returns a sql.DB connection pool for future use.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// as connections to the database are only established lazily when they're needed,
	// and because sql.Open() does not verify that the database is actually accessible,
	// we need to explicitly ping the database to check that it's reachable and the DSN is valid.
	// If we don't do this, then any issues with the database connection won't be discovered
	// until we try to execute a query against the database later on in our code,
	// which could lead to unexpected errors at runtime.
	// By pinging the database immediately after opening the connection pool, we can catch any connection issues early on and handle them appropriately
	// (e.g., by logging an error message and exiting the application).
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
