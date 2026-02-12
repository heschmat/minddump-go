package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type config struct {
	addr      string
	staticDir string
}

// a struct to hold the application-wide dependencies
type application struct {
	logger *slog.Logger
}

func main() {
	// flags & env variables --------------------
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	// logger -----------------------------------
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // for development
		// AddSource: true,
	}))

	// application dependencies -----------------
	app := &application{
		logger: logger,
	}

	// servemux ---------------------------------
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/")) // &{./ui/static/}
	// N.B. actual file access happens per request.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.getSnippetByID)

	// server -----------------------------------
	logger.Info("starting server...", "addr", cfg.addr)

	// params: the TCP network address, the servemux
	err := http.ListenAndServe(cfg.addr, mux)
	// N.B. Any error returned by http.ListenAndServe() is always non-nil.
	// log.Fatal(err)
	logger.Error(err.Error())
	os.Exit(1)
}
