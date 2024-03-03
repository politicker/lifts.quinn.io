package web

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/politicker/lifts.quinn.io/internal/db"
	"go.uber.org/zap"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed templates/index.html
var indexHTML embed.FS

type LiftData struct {
	Name    string
	BestSet string
	History []float64
}

type TemplateData struct {
	LiftData string
}

type web struct {
	port       int
	logger     *zap.Logger
	queries    *db.Queries
	httpClient *http.Client
}

func NewWeb(ctx context.Context, logger *zap.Logger, database *sql.DB, port int) *web {
	return &web{
		port:       port,
		logger:     logger,
		queries:    db.New(database),
		httpClient: &http.Client{},
	}
}

func (s *web) Start() error {
	http.HandleFunc("/", s.indexHandler)
	fs := http.FileServer(http.FS(staticFiles))
	http.Handle("/static/", fs)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

func (s *web) indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(indexHTML, "templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	liftData, err := s.getLiftData(r.Context())
	if err != nil {
	}

	jsonData, err := json.Marshal(liftData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, TemplateData{LiftData: string(jsonData)})
	if err != nil {
	}
}

func (s *web) getLiftData(ctx context.Context) ([]LiftData, error) {
	ohpHistory, err := s.queries.Get1RMHistory(ctx, "overhead press (barbell)")
	if err != nil {
		s.logger.Error("failed to get ohp history")
	}

	squatHistory, err := s.queries.Get1RMHistory(ctx, "squat (barbell)")
	if err != nil {
		s.logger.Error("failed to get squat history")
	}

	benchHistory, err := s.queries.Get1RMHistory(ctx, "bench press (barbell)")
	if err != nil {
		s.logger.Error("failed to get bench history")
	}

	deadliftHistory, err := s.queries.Get1RMHistory(ctx, "deadlift (barbell)")
	if err != nil {
		s.logger.Error("failed to get deadlift history")
	}

	bestOHP, err := s.queries.GetBestSet(ctx, "overhead press (barbell)")
	if err != nil {
		s.logger.Error("failed to get best ohp")
	}

	bestSquat, err := s.queries.GetBestSet(ctx, "squat (barbell)")
	if err != nil {
		s.logger.Error("failed to get best squat")
	}

	bestBench, err := s.queries.GetBestSet(ctx, "bench press (barbell)")
	if err != nil {
		s.logger.Error("failed to get best bench")
	}

	bestDeadlift, err := s.queries.GetBestSet(ctx, "deadlift (barbell)")
	if err != nil {
		s.logger.Error("failed to get best deadlift")
	}

	return []LiftData{{
		Name:    "OHP",
		BestSet: bestOHP.(string),
		History: ohpHistory,
	}, {
		Name:    "Squat",
		BestSet: bestSquat.(string),
		History: squatHistory,
	}, {
		Name:    "Bench",
		BestSet: bestBench.(string),
		History: benchHistory,
	}, {
		Name:    "Deadlift",
		BestSet: bestDeadlift.(string),
		History: deadliftHistory,
	}}, nil
}
