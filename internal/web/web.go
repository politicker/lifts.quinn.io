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

const (
	deadlift = "deadlift (barbell)"
	bench    = "bench press (barbell)"
	squat    = "squat (barbell)"
	ohp      = "overhead press (barbell)"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(liftData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, TemplateData{LiftData: string(jsonData)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *web) getLiftData(ctx context.Context) ([]LiftData, error) {
	ohpHistory, err := s.queries.Get1RMHistory(ctx, ohp)
	if err != nil {
		s.logger.Error("failed to get ohp history")
		return nil, err
	}

	squatHistory, err := s.queries.Get1RMHistory(ctx, squat)
	if err != nil {
		s.logger.Error("failed to get squat history")
		return nil, err
	}

	benchHistory, err := s.queries.Get1RMHistory(ctx, bench)
	if err != nil {
		s.logger.Error("failed to get bench history")
		return nil, err
	}

	deadliftHistory, err := s.queries.Get1RMHistory(ctx, deadlift)
	if err != nil {
		s.logger.Error("failed to get deadlift history")
		return nil, err
	}

	bestOHP, err := s.queries.GetBestSet(ctx, ohp)
	if err != nil {
		s.logger.Error("failed to get best ohp")
		return nil, err
	}

	bestSquat, err := s.queries.GetBestSet(ctx, squat)
	if err != nil {
		s.logger.Error("failed to get best squat")
		return nil, err
	}

	bestBench, err := s.queries.GetBestSet(ctx, bench)
	if err != nil {
		s.logger.Error("failed to get best bench")
		return nil, err
	}

	bestDeadlift, err := s.queries.GetBestSet(ctx, deadlift)
	if err != nil {
		s.logger.Error("failed to get best deadlift")
		return nil, err
	}

	return []LiftData{{
		Name:    "OHP",
		BestSet: bestOHP.String,
		History: ohpHistory,
	}, {
		Name:    "Squat",
		BestSet: bestSquat.String,
		History: squatHistory,
	}, {
		Name:    "Bench",
		BestSet: bestBench.String,
		History: benchHistory,
	}, {
		Name:    "Deadlift",
		BestSet: bestDeadlift.String,
		History: deadliftHistory,
	}}, nil
}
