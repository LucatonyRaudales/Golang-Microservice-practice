package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nicholasjackson/env"
)
var binAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

func main(){
	env.Parse()

	l := log.New(os.Stdout, "products-api", log.LstdFlags)

	// create the handlers
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodBye()

	// create a new serve mux and register the handler
	sm := http.NewServeMux()
	sm.HandleFunc("/", hh)
	sm.HandleFunc("/goodbye", gh)

	// create a new server
	s := http.Server{
		Addr: *binAddress,
		Handler: sm,
		ErrorLog: l,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 120 * time.Second,
	}

	// start the server
	go func ()  {
		l.Println("Starting server on port 9090")
		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm ir interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// block until a signal is received
	sig := <-c
	log.Println("Got signl:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30 * time.Second)
	s.Shutdown(ctx)
}