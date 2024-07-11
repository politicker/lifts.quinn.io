package domain

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/politicker/lifts.quinn.io/internal/db"
	"go.uber.org/zap"
)

type Set struct {
	LoggedAt     string `csv:"Date"`
	WorkoutName  string `csv:"Workout Name"`
	Duration     string `csv:"Duration"`
	ExerciseName string `csv:"Exercise Name"`
	SetOrder     string `csv:"Set Order"`
	Weight       string `csv:"Weight"`
	Reps         string `csv:"Reps"`
	Distance     string `csv:"Distance"`
	Seconds      string `csv:"Seconds"`
	Notes        string `csv:"Notes"`
	WorkoutNotes string `csv:"Workout Notes"`
	RPE          string `csv:"RPE"`
}

type importer struct {
	logger  *zap.Logger
	queries *db.Queries
}

func NewImporter(logger *zap.Logger, queries *db.Queries) *importer {
	return &importer{
		logger:  logger,
		queries: queries,
	}
}

func (i *importer) Run(ctx context.Context, reader io.Reader) error {
	sets := []Set{}
	if err := gocsv.Unmarshal(reader, &sets); err != nil {
		return err
	}

	if len(sets) == 0 {
		return errors.New("lifts CSV file contains no data")
	}

	for _, set := range sets {
		if set.SetOrder == "Rest Timer" {
			i.logger.Debug("skipping rest timer set")
			continue
		}

		weight, err := strconv.ParseFloat(set.Weight, 64)
		if err != nil {
			i.logger.Error("failed to parse weight", zap.Error(err))
			continue
		}

		reps, err := strconv.ParseFloat(set.Reps, 64)
		if err != nil {
			i.logger.Error("failed to parse reps", zap.Error(err))
			continue
		}

		seconds, err := strconv.ParseFloat(set.Seconds, 64)
		if err != nil {
			i.logger.Error("failed to parse seconds", zap.Error(err))
			continue
		}

		setOrder, err := strconv.ParseInt(set.SetOrder, 10, 64)
		if err != nil {
			i.logger.Error("failed to parse set order", zap.Error(err), zap.String("set_order", set.SetOrder))
			continue
		}

		loggedAt, err := time.Parse("2006-01-02 15:04:05", set.LoggedAt)
		if err != nil {
			i.logger.Error("failed to parse loggedAt", zap.Error(err))
			continue
		}

		err = i.queries.CreateLiftSetLog(ctx, db.CreateLiftSetLogParams{
			WorkoutName:  set.WorkoutName,
			ExerciseName: set.ExerciseName,
			Weight:       weight,
			Reps:         reps,
			Seconds:      seconds,
			SetOrder:     int32(setOrder),
			Notes:        sql.NullString{String: set.Notes, Valid: set.Notes != ""},
			WorkoutNotes: sql.NullString{String: set.WorkoutNotes, Valid: set.WorkoutNotes != ""},
			LoggedAt:     loggedAt,
		})
		if err != nil {
			i.logger.Error("failed to create lift set log", zap.Error(err), zap.String("workout_name", set.WorkoutName), zap.String("exercise_name", set.ExerciseName))
			continue
		}
	}

	return nil
}
