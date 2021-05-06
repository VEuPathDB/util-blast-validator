package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"blast-validator/v1/internal/server"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	server.RegisterEndpoint(router)

	i1 := make(chan os.Signal, 1)
	i2 := make(chan error, 1)

	signal.Notify(i1, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Starting server on port 80.")
		i2 <- http.ListenAndServe("0.0.0.0:80", router)
	}()

	select {
	case <- i1:
		log.Println("Shutdown signal received.")
		os.Exit(0)
	case e := <- i2:
		log.Println("Server crashed: ", e)
		os.Exit(0)
	}
}
