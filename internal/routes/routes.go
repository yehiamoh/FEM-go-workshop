package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/yehiamoh/go-fem-workshop/internal/app"
)

func SetUpRoutes(app *app.Application) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/health", app.HealthCheck)
	registerWorkoutRoutes(router, app)
	registerUserRoutes(router, app)
	registerTokenRoutes(router, app)
	return router
}

func registerWorkoutRoutes(router *chi.Mux, app *app.Application) {
	router.Route("/workouts", func(r chi.Router) {
		r.Post("/", app.WorkoutHandler.HandleCreateWorkout)
		r.Get("/{id}", app.WorkoutHandler.HandleGetWorkOutByID)
		r.Put("/{id}", app.WorkoutHandler.HandleUpdateWorkoutByID)
		r.Delete("/{id}", app.WorkoutHandler.HandleDeleteWorkout)
	})
}
func registerUserRoutes(router *chi.Mux, app *app.Application) {
	router.Route("/users", func(r chi.Router) {
		r.Post("/", app.UserHandler.HandleRegisterUser)
	})
}
func registerTokenRoutes(router *chi.Mux, app *app.Application) {
	router.Route("/token/authentication", func(r chi.Router) {
		r.Post("/", app.TokenHandler.HandleCreateToken)
	})
}
