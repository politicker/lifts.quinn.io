// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"database/sql"
	"time"
)

type LiftSetLog struct {
	WorkoutName     string
	WorkoutDuration string
	ExerciseName    string
	SetOrder        int32
	Weight          float64
	Reps            float64
	Distance        float64
	Seconds         float64
	Notes           sql.NullString
	WorkoutNotes    sql.NullString
	Rpe             sql.NullString
	LoggedAt        time.Time
	ImportedAt      time.Time
}
