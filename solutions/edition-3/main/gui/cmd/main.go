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
	settings := flag.String("settings", "", "GUI settings JSON path (default: cr/settings.json under the working directory)")
	save := flag.String("save", "", "recovery archive (default: save.json)")
	logMode := flag.String("log-mode", "", "archive retention: full or dialogue_only")
	pause := flag.Bool("pause-on-start", false, "pause before the first model request")
	speaking := flag.Bool("speaking", false, "enable Agent speech requests")
	flag.Parse()
	config, err := ensemble.ConfigFromEnv(os.Getenv)
	if err != nil {
		log.Fatal(err)
	}
	if config.SavePath == "" {
		config.SavePath = "save.json"
	}
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "save":
			config.SavePath = *save
		case "log-mode":
			config.LogMode = ensemble.LogMode(*logMode)
		case "pause-on-start":
			config.PauseOnStart = *pause
		case "speaking":
			config.Speaking = *speaking
		}
	})
	root := ensemble.New(log.Default())
	var trace *ensemble.DiagnosticLog
	ready, normalExit := false, false
	defer func() {
		if err := root.Shutdown(); err != nil {
			log.Printf("shutdown/save: %v", err)
			return
		}
		if ready && normalExit && trace != nil {
			if err := trace.Retain(); err != nil {
				log.Printf("diagnostic retention: %v", err)
			}
		}
	}()
	agent, err := root.NewAgent(config)
	if err != nil {
		log.Fatal(err)
	}
	if r := agent.Recovery(); r.Path != "" {
		log.Printf("loaded %s through sequence %d (backup=%t)", r.Path, r.AsOf, r.Backup)
	}
	app, err := gui.NewServerWithSettings(root, *assets, *settings)
	if err != nil {
		log.Print(err)
		return
	}
	defer app.Close()
	app.DefaultAgent = agent.ID()
	trace, err = ensemble.OpenDiagnosticLog("gui.log", func() int { return *agent.Config().DiagnosticMaxRecords })
	if err != nil {
		log.Print(err)
		return
	}
	defer trace.Close()
	app.Trace = trace
	if agent.Recovery().Path == "" {
		if err = agent.Save(); err != nil {
			log.Printf("startup save: %v", err)
			return
		}
	}
	if err = trace.Retain(); err != nil {
		log.Print(err)
		return
	}
	ready = true
	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	terminalDone := app.Terminal(ctx, os.Stdin, os.Stdout)
	server := &http.Server{Addr: addr, Handler: app.Handler(), ReadHeaderTimeout: 5 * time.Second}
	drained := make(chan bool, 1)
	go func() {
		completed := false
		select {
		case <-ctx.Done():
		case <-terminalDone:
			completed = ctx.Err() == nil
		}
		drained <- completed
		stop()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdown)
	}()
	log.Printf("Ensemble GUI http://%s · assets %s · terminal EOF saves and exits", addr, app.Assets)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Print(err)
	} else {
		normalExit = <-drained
	}
}
