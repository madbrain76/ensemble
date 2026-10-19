package main

import (
	"context"
	"ensemble"
	"ensemble/gui"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
func run() error {
	c, err := ensemble.ConfigFromEnv(os.Getenv)
	if err != nil {
		return err
	}
	app := ensemble.New(log.Default())
	defer app.Shutdown()
	a, err := app.NewAgent(c)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	server := &http.Server{Addr: "127.0.0.1:8080", Handler: gui.New(a)}
	ended := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = app.Shutdown()
			_ = server.Close()
		case <-ended:
		}
	}()
	defer close(ended)
	log.Print("GUI listening on http://127.0.0.1:8080")
	err = server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
