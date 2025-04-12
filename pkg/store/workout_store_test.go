package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setuptestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5434 sslmode=disable")
	if err != nil {
		t.Fatalf("opening test db :%v", err)
	}
	err = Migrate(db, "../../migrations")
	if err != nil {
		t.Fatalf("migration test db error:%v", err)
	}
	_, err = db.Exec(`TRUNCATE workouts,workout_entries CASCADE`)
	if err != nil {
		t.Fatalf("turncating tables :%v", err)
	}
	return db
}
func TestCreateWorkout(t *testing.T) {
	db := setuptestDB(t)
	defer db.Close()
	store := NewPostgresWorkoutStore(db)
	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:           "push day",
				Description:     "upper body day",
				DurationMinutes: 60,
				CaloriesBurned:  200,
				Entries: []WorkoutEntry{
					{
						ExerciesName: "Bench press",
						Sets:         3,
						Reps:         ptrInt(100),
						Weight:       ptrFloat(135.5),
						Notes:        "Warm up properly",
						OrderIndex:   1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "not valid workout",
			workout: &Workout{
				Title:           "Full Body",
				Description:     "full body day",
				DurationMinutes: 80,
				CaloriesBurned:  400,
				Entries: []WorkoutEntry{
					{
						ExerciesName: "Plank",
						Sets:         3,
						Reps:         ptrInt(60),
						//Weight:       ptrFloat(200.5),
						Notes:      "keep form",
						OrderIndex: 1,
					},
					{
						ExerciesName:    "Squats",
						Sets:            7,
						Reps:            ptrInt(60),
						DurationSeconds: ptrInt(60),
						Weight:          ptrFloat(200.5),
						Notes:           "full depth",
						OrderIndex:      2,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(testCase.workout)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testCase.workout.Title, createdWorkout.Title)
			assert.Equal(t, testCase.workout.Description, createdWorkout.Description)
			assert.Equal(t, testCase.workout.DurationMinutes, createdWorkout.DurationMinutes)
			retrived, err := store.GetWorkoutByID(int(createdWorkout.ID))
			assert.Equal(t, createdWorkout.ID, retrived.ID)
			assert.Equal(t, len(testCase.workout.Entries), len(retrived.Entries))
			for i := range retrived.Entries {
				assert.Equal(t, testCase.workout.Entries[i].ExerciesName, retrived.Entries[i].ExerciesName)
				assert.Equal(t, testCase.workout.Entries[i].Sets, retrived.Entries[i].Sets)
				assert.Equal(t, testCase.workout.Entries[i].OrderIndex, retrived.Entries[i].OrderIndex)
			}

		})

	}
}
func ptrInt(i int) *int {
	return &i
}
func ptrFloat(i float64) *float64 {
	return &i
}
