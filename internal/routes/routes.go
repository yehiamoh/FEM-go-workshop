package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/yehiamoh/go-fem-workshop/internal/app"
)

func SetUpRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", app.HealthCheck)

	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
	r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkOutByID)

	return r
}
