package web

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/politicker/lifts.quinn.io/internal/db"
	"github.com/politicker/lifts.quinn.io/internal/domain"
	"go.uber.org/zap"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed templates/index.html
var indexHTML embed.FS

type TemplateData struct {
	Title string
	Lifts string
}

type web struct {
	port       int
	title      string
	logger     *zap.Logger
	queries    *db.Queries
	httpClient *http.Client
}

func NewWeb(ctx context.Context, logger *zap.Logger, database *sql.DB, port int, title string) *web {
	return &web{
		port:       port,
		title:      title,
		logger:     logger,
		queries:    db.New(database),
		httpClient: &http.Client{},
	}
}

func (s *web) Start() error {
	http.HandleFunc("/", s.indexHandler)
	http.HandleFunc("/upload-lifts", s.uploadLiftsHandler)

	fs := http.FileServer(http.FS(staticFiles))
	http.Handle("/static/", fs)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

func (s *web) uploadLiftsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		s.logger.Info("Method not allowed", zap.String("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ensure the content type is correct
	if r.Header.Get("Content-Type") != "text/csv" {
		s.logger.Info("Unsupported media type", zap.String("content-type", r.Header.Get("Content-Type")))
		http.Error(w, "Unsupported media type. Please upload a CSV file.", http.StatusUnsupportedMediaType)
		return
	}

	userSecret := r.Header.Get("X-Upload-Secret")
	uploadSecret := os.Getenv("UPLOAD_SECRET")
	if subtle.ConstantTimeCompare([]byte(userSecret), []byte(uploadSecret)) != 1 {
		s.logger.Info("Unauthorized", zap.String("secret", userSecret))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	s.logger.Info("Received lift data upload")
	importer := domain.NewImporter(s.logger, s.queries)
	err := importer.Run(r.Context(), r.Body)
	if err != nil {
		s.logger.Error("Error importing lift data", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.logger.Info("Lift data imported successfully")
	w.WriteHeader(http.StatusOK)
}

func (s *web) indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(indexHTML, "templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lr := domain.NewLiftsRepository(s.logger, s.queries)
	liftData, err := lr.GetLifts(r.Context())
	if err != nil {
		// TODO: Fix this for real
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

	jsonData, err := json.Marshal(liftData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, TemplateData{Title: s.title, Lifts: string(jsonData)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
