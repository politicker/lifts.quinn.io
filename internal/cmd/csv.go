package cmd

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/politicker/lifts.quinn.io/internal/cmdutil"
	"github.com/politicker/lifts.quinn.io/internal/db"
	"github.com/spf13/cobra"
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

func ExtractCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract",
		Args:  cobra.ExactArgs(0),
		Short: "Import lift history from a CSV file",
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath, err := cmd.Flags().GetString("file")
			if err != nil {
				return err
			}

			if filePath == "" {
				return errors.New("file is a required flag")
			}

			liftCSVFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
			if err != nil {
				panic(err)
			}
			defer liftCSVFile.Close()

			fileInfo, err := liftCSVFile.Stat()
			if err != nil {
				panic(err)
			}
			if fileInfo.Size() == 0 {
				return errors.New("lifts CSV file is empty")
			}
			sets := []Set{}
			if err = gocsv.UnmarshalFile(liftCSVFile, &sets); err != nil {
				panic(err)
			}

			logger := cmdutil.NewLogger("import")
			defer func() { _ = logger.Sync() }()

			database, err := cmdutil.NewDBConnection(ctx)
			if err != nil {
				return err
			}
			defer database.Close()

			queries := db.New(database)

			for _, set := range sets {
				weight, err := strconv.ParseFloat(set.Weight, 64)
				if err != nil {
					logger.Error("failed to parse weight", zap.Error(err))
					continue
				}

				reps, err := strconv.ParseFloat(set.Reps, 64)
				if err != nil {
					logger.Error("failed to parse reps", zap.Error(err))
					continue
				}

				seconds, err := strconv.ParseFloat(set.Seconds, 64)
				if err != nil {
					logger.Error("failed to parse seconds", zap.Error(err))
					continue
				}

				loggedAt, err := time.Parse("2006-01-02 15:04:05", set.LoggedAt)
				if err != nil {
					logger.Error("failed to parse loggedAt", zap.Error(err))
					continue
				}

				err = queries.CreateLiftSetLog(ctx, db.CreateLiftSetLogParams{
					WorkoutName:  set.WorkoutName,
					ExerciseName: set.ExerciseName,
					Weight:       weight,
					Reps:         reps,
					Seconds:      seconds,
					Notes:        sql.NullString{String: set.Notes, Valid: set.Notes != ""},
					WorkoutNotes: sql.NullString{String: set.WorkoutNotes, Valid: set.WorkoutNotes != ""},
					LoggedAt:     sql.NullTime{Time: loggedAt, Valid: true},
				})
				if err != nil {
					logger.Error("failed to create lift set log", zap.Error(err), zap.String("workout_name", set.WorkoutName), zap.String("exercise_name", set.ExerciseName))
					continue
				}
			}

			logger.Info("imported workout history")
			return nil
		},
	}

	cmd.Flags().StringP("file", "f", "", "The CSV file to import")
	return cmd
}
