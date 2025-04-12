package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yehiamoh/go-fem-workshop/migrations"
	"github.com/yehiamoh/go-fem-workshop/pkg/api"
	"github.com/yehiamoh/go-fem-workshop/pkg/middleware"
	"github.com/yehiamoh/go-fem-workshop/pkg/store"
)

type Application struct {
	Logger           *log.Logger
	WorkoutHandler   *api.WorkoutHandler
	UserHandler      *api.UserHandler
	TokenHandler     *api.TokenHandler
	MiddlwareHandler *middleware.UserMiddleware
	DB               *sql.DB
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

	userStore := store.NewPostgresUserStore(pgDB)
	userHandler := api.NewUserHandler(userStore, logger)

	tokenSotre := store.NewPostgresTokenStore(pgDB)
	tokenHandler := api.NewTokenHandler(tokenSotre, userStore, logger)

	middlewareHandler := &middleware.UserMiddleware{
		UserStore: userStore,
	}

	app := &Application{
		Logger:           logger,
		WorkoutHandler:   workoutHandler,
		UserHandler:      userHandler,
		TokenHandler:     tokenHandler,
		MiddlwareHandler: middlewareHandler,
		DB:               pgDB,
	}
	return app, nil
}
func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
