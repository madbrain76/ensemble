// gui hosts an optional browser and terminal view through Ensemble's public API.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ensemble"
	"ensemble/gui"
)

func main() {
	port := flag.Int("port", 8088, "loopback HTTP port")
	assets := flag.String("gui-dir", "web/gui", "GUI asset directory (absolute paths work outside cwd)")
	pause := flag.Bool("pause-on-start", false, "pause before the first model request")
	speaking := flag.Bool("speaking", false, "enable Agent speech requests")
	flag.Parse()
	config, err := ensemble.ConfigFromEnv(os.Getenv)
	if err != nil {
		log.Fatal(err)
	}
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "pause-on-start":
			config.PauseOnStart = *pause
		case "speaking":
			config.Speaking = *speaking
		}
	})
	root := ensemble.New(log.Default())
	defer root.Shutdown()
	agent, err := root.NewAgent(config)
	if err != nil {
		log.Fatal(err)
	}
	app, err := gui.NewServer(root, *assets)
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()
	app.DefaultAgent = agent.ID()
	trace, err := os.OpenFile("gui.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer trace.Close()
	app.Trace = trace
	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	app.Terminal(ctx, os.Stdin, os.Stdout)
	server := &http.Server{Addr: addr, Handler: app.Handler(), ReadHeaderTimeout: 5 * time.Second}
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdown)
	}()
	log.Printf("Ensemble GUI http://%s · assets %s · terminal EOF detaches input", addr, app.Assets)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Print(err)
	}
}
