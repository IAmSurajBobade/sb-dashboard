package main

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback to 8080 if not set
	}
	logLevel := os.Getenv("LOG_LEVEL")
	var slogLevel slog.Level
	switch logLevel {
	case "DEBUG":
		slogLevel = slog.LevelDebug
	case "INFO":
		slogLevel = slog.LevelInfo
	case "WARN":
		slogLevel = slog.LevelWarn
	case "ERROR":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo // fallback to INFO if not set
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slogLevel,
		AddSource: true,
	})).With(
		"service", "sb-dashboard",
	)
	router := http.NewServeMux()
	router.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "OK"}`))
	}))
	router.Handle("/info", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ist, err := time.LoadLocation("Asia/Kolkata")
		if err != nil {
			panic(err)
		}
		bdate := time.Date(1995, 5, 28, 0, 0, 0, 0, ist)
		daysSince := time.Since(bdate).Hours() / 24
		logger.Debug("/info called", slog.String("daysSince", fmt.Sprintf("%.0f", daysSince)))
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "daysSince": "` + fmt.Sprintf("%.0f", daysSince) + `"}`))
	}))

	// Then use it for your HTTP server
	server := &http.Server{
		Addr:     net.JoinHostPort("", port),
		Handler:  router,
		ErrorLog: slog.NewLogLogger(logger.Handler(), slogLevel),
	}
	// And start it
	logger.Info("starting server", "addr", server.Addr)
	server.ListenAndServe()
}
