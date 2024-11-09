package main

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
)

// APIResponse represents our API data structure
type APIResponse struct {
	DaysSince string `json:"daysSince"`
}

// PageData represents data we'll pass to our HTML template
type PageData struct {
	Endpoints []EndpointInfo
}

type EndpointInfo struct {
	Path        string
	Description string
	Example     string
}

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

	// Add template handling
	tmpl, err := template.New("index").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>SB Dashboard API</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 40px auto;
            max-width: 800px;
            padding: 0 20px;
        }
        .endpoint {
            border: 1px solid #ddd;
            margin: 20px 0;
            padding: 20px;
            border-radius: 5px;
        }
        .endpoint h3 {
            margin-top: 0;
            color: #333;
        }
        .example {
            background-color: #f5f5f5;
            padding: 10px;
            border-radius: 3px;
            font-family: monospace;
        }
    </style>
</head>
<body>
    <h1>SB Dashboard API Documentation</h1>
    <div class="endpoints">
        {{range .Endpoints}}
        <div class="endpoint">
            <h3>{{.Path}}</h3>
            <p>{{.Description}}</p>
            <div class="example">
                <strong>Example Response:</strong><br>
                {{.Example}}
            </div>
        </div>
        {{end}}
    </div>
</body>
</html>
`)
	if err != nil {
		logger.Error("failed to parse template", "error", err)
		panic(err)
	}

	// Handler function to check if request is from browser
	isBrowser := func(r *http.Request) bool {
		userAgent := r.Header.Get("User-Agent")
		return strings.Contains(userAgent, "Mozilla") ||
			strings.Contains(userAgent, "Chrome") ||
			strings.Contains(userAgent, "Safari") ||
			strings.Contains(userAgent, "Edge") ||
			strings.Contains(userAgent, "Firefox")
	}

	// Root handler
	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		pageData := PageData{
			Endpoints: []EndpointInfo{
				{
					Path:        "/health",
					Description: "Health check endpoint",
					Example:     `{"status": "OK"}`,
				},
				{
					Path:        "/info",
					Description: "Returns the number of days since a specific date",
					Example:     `{"daysSince": "10000"}`,
				},
			},
		}

		w.Header().Set("Content-Type", "text/html")
		err := tmpl.Execute(w, pageData)
		if err != nil {
			logger.Error("failed to execute template", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}))

	// Update health endpoint to handle both browser and API requests
	router.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{"status": "OK"}`

		if isBrowser(r) {
			w.Header().Set("Content-Type", "text/html")
			tmpl.Execute(w, PageData{
				Endpoints: []EndpointInfo{
					{
						Path:        "/health",
						Description: "Health check endpoint",
						Example:     response,
					},
				},
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))

	// Update info endpoint to handle both browser and API requests
	router.Handle("/info", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ist, err := time.LoadLocation("Asia/Kolkata")
		if err != nil {
			logger.Error("failed to load timezone", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		bdate := time.Date(1995, 5, 28, 0, 0, 0, 0, ist)
		daysSince := time.Since(bdate).Hours() / 24
		response := fmt.Sprintf(`{"daysSince": "%.0f"}`, daysSince)

		if isBrowser(r) {
			w.Header().Set("Content-Type", "text/html")
			tmpl.Execute(w, PageData{
				Endpoints: []EndpointInfo{
					{
						Path:        "/info",
						Description: "Returns the number of days since May 28, 1995",
						Example:     response,
					},
				},
			})
			return
		}

		logger.Debug("/info called", slog.String("daysSince", fmt.Sprintf("%.0f", daysSince)))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))

	// Your existing server setup
	server := &http.Server{
		Addr:     net.JoinHostPort("", port),
		Handler:  router,
		ErrorLog: slog.NewLogLogger(logger.Handler(), slogLevel),
	}

	logger.Info("starting server", "addr", server.Addr)
	server.ListenAndServe()
}
