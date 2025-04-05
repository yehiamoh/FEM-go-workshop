package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yehiamoh/go-fem-workshop/internal/api"
	"github.com/yehiamoh/go-fem-workshop/internal/store"
	"github.com/yehiamoh/go-fem-workshop/migrations"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFs(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		DB:             pgDB,
	}
	return app, nil
}
func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
