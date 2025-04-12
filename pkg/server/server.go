package server

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/yehiamoh/go-fem-workshop/pkg/app"
	"github.com/yehiamoh/go-fem-workshop/pkg/routes"
)

func NewHttpServer(port int, router *chi.Mux) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", port),
		IdleTimeout:  time.Minute,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
func Run() error {
	var port int

	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse()

	myApp, err := app.NewApplication()
	if err != nil {
		return err
	}
	defer myApp.DB.Close()
	router := routes.SetUpRoutes(myApp)
	myApp.Logger.Printf("App running on port : %v", port)

	server := NewHttpServer(port, router)

	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
