package domain

import (
	"context"

	"github.com/politicker/lifts.quinn.io/internal/db"
	"go.uber.org/zap"
)

type liftsRepository struct {
	logger  *zap.Logger
	queries *db.Queries
}

type Lift struct {
	Name    string
	BestSet string
	History []float64
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

func (lr *liftsRepository) GetLifts(ctx context.Context) ([]Lift, error) {
	ohpHistory, err := lr.queries.Get1RMHistory(ctx, ohp)
	if err != nil {
		lr.logger.Error("failed to get ohp history")
		return nil, err
	}

	squatHistory, err := lr.queries.Get1RMHistory(ctx, squat)
	if err != nil {
		lr.logger.Error("failed to get squat history")
		return nil, err
	}

	benchHistory, err := lr.queries.Get1RMHistory(ctx, bench)
	if err != nil {
		lr.logger.Error("failed to get bench history")
		return nil, err
	}

	deadliftHistory, err := lr.queries.Get1RMHistory(ctx, deadlift)
	if err != nil {
		lr.logger.Error("failed to get deadlift history")
		return nil, err
	}

	bestOHP, err := lr.queries.GetBestSet(ctx, ohp)
	if err != nil {
		lr.logger.Error("failed to get best ohp")
		return nil, err
	}

	bestSquat, err := lr.queries.GetBestSet(ctx, squat)
	if err != nil {
		lr.logger.Error("failed to get best squat")
		return nil, err
	}

	bestBench, err := lr.queries.GetBestSet(ctx, bench)
	if err != nil {
		lr.logger.Error("failed to get best bench")
		return nil, err
	}

	bestDeadlift, err := lr.queries.GetBestSet(ctx, deadlift)
	if err != nil {
		lr.logger.Error("failed to get best deadlift")
		return nil, err
	}

	return []Lift{{
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
