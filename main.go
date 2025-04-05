package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/yehiamoh/go-fem-workshop/internal/app"
	"github.com/yehiamoh/go-fem-workshop/internal/routes"
)

func main() {

	var port int

	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse()

	/*
		The flag package in Go is used for parsing command-line arguments. In your code, it's being used to define and parse a port number flag:
		# Use default port (8080)
			go run main.go
		# Use custom port
			go run main.go -port 3000
	*/

	myApp, err := app.NewApplication()
	if err != nil {
		panic(err.Error())
	}

	defer myApp.DB.Close()

	router := routes.SetUpRoutes(myApp)

	myApp.Logger.Printf("Running Our App on port %d", port)

	server := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", port),
		IdleTimeout:  time.Minute,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		myApp.Logger.Fatal(err.Error())
	}

}
