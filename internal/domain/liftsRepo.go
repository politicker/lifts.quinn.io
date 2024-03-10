package domain

import (
	"context"
	"time"

	"github.com/politicker/lifts.quinn.io/internal/db"
	"go.uber.org/zap"
)

type liftsRepository struct {
	logger  *zap.Logger
	queries *db.Queries
}

type LiftHistory struct {
	Name    string
	BestSet Lift
	History []Lift
}

type Lift struct {
	SetText string  // "200x5"
	RepMax  float64 // 1rm
	Date    time.Time
}

const (
	deadlift = "deadlift (barbell)"
	bench    = "bench press (barbell)"
	squat    = "squat (barbell)"
	ohp      = "overhead press (barbell)"
)

func NewLiftsRepository(logger *zap.Logger, queries *db.Queries) *liftsRepository {
	return &liftsRepository{
		logger:  logger,
		queries: queries,
	}
}

func (lr *liftsRepository) GetLifts(ctx context.Context) ([]LiftHistory, error) {
	bestOHP, err := lr.getBestSet(ctx, ohp)
	if err != nil {
		return nil, err
	}

	bestSquat, err := lr.getBestSet(ctx, squat)
	if err != nil {
		return nil, err
	}

	bestBench, err := lr.getBestSet(ctx, bench)
	if err != nil {
		return nil, err
	}

	bestDeadlift, err := lr.getBestSet(ctx, deadlift)
	if err != nil {
		return nil, err
	}

	ohpLiftHistory, err := lr.getLiftHistory(ctx, ohp)
	if err != nil {
		return nil, err
	}

	squatHistory, err := lr.getLiftHistory(ctx, squat)
	if err != nil {
		return nil, err
	}

	benchHistory, err := lr.getLiftHistory(ctx, bench)
	if err != nil {
		return nil, err
	}

	deadliftHistory, err := lr.getLiftHistory(ctx, deadlift)
	if err != nil {
		return nil, err
	}

	return []LiftHistory{{
		Name:    "OHP",
		BestSet: bestOHP,
		History: ohpLiftHistory,
	}, {
		Name:    "Squat",
		BestSet: bestSquat,
		History: squatHistory,
	}, {
		Name:    "Bench",
		BestSet: bestBench,
		History: benchHistory,
	}, {
		Name:    "Deadlift",
		BestSet: bestDeadlift,
		History: deadliftHistory,
	}}, nil
}

func (lr *liftsRepository) getBestSet(ctx context.Context, liftType string) (Lift, error) {
	bestSet, err := lr.queries.GetBestSet(ctx, liftType)
	if err != nil {
		lr.logger.Error("failed to get best "+liftType, zap.Error(err))
		return Lift{}, err
	}

	return Lift{
		SetText: bestSet.SetText.String,
		RepMax:  bestSet.Estimated1rm,
		Date:    bestSet.LoggedAt,
	}, nil
}

func (lr *liftsRepository) getLiftHistory(ctx context.Context, liftType string) ([]Lift, error) {
	liftHistory, err := lr.queries.Get1RMHistory(ctx, liftType)
	if err != nil {
		lr.logger.Error("failed to get "+liftType+" history", zap.Error(err))
		return nil, err
	}

	var history []Lift
	for _, h := range liftHistory {
		history = append(history, Lift{
			SetText: h.SetText.String,
			RepMax:  h.Estimated1rm,
			Date:    h.LoggedAt,
		})
	}

	return history, nil
}
