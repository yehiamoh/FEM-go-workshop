package store

import (
	"database/sql"
)

type Workout struct {
	ID              int            `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	ID              int      `json:"id"`
	ExerciesName    string   `json:"exercies_name"`
	Sets            int      `json:"int"`
	Reps            *int     `json:"reps"` //because we will explcilty chech the reps
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight"`
	Notes           string   `json:"string"`
	OrderIndex      int      `json:"order_index"`
}
type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkoutByID(id int) (*Workout, error)
}

func (pg *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query :=
		`INSERT INTO worouts(title,description,duration_minutes,calories_burned)
	VALUES($1,$2,$3,$4)
	RETURNING id
	`
	err = tx.QueryRow(query, workout.ID, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)
	if err != nil {
		return nil, err
	}
	for _, entry := range workout.Entries {
		query := `
		INSERT INTO workout_entries(workout_id,exercise_name,sets,reps,duration_seconds,weight,notes,order_index)
		VALUES($1,$2,$3,$4,$5,$6)
		RETURNING id
		`
		err := tx.QueryRow(query, workout.ID, entry.ExerciesName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return workout, nil
}
