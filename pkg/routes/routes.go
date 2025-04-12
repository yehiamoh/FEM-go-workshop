package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/yehiamoh/go-fem-workshop/pkg/app"
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
	router.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkOutByID)
	router.Route("/workouts", func(r chi.Router) {
		r.Use(app.MiddlwareHandler.Authenticate)
		r.Post("/", app.MiddlwareHandler.RequireUser(app.WorkoutHandler.HandleCreateWorkout))
		r.Put("/{id}", app.MiddlwareHandler.RequireUser(app.WorkoutHandler.HandleUpdateWorkoutByID))
		r.Delete("/{id}", app.MiddlwareHandler.RequireUser(app.WorkoutHandler.HandleDeleteWorkout))
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
