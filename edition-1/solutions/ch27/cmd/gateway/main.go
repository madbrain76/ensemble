// Command gateway runs the Ensemble gateway: the edge service that sits at
// the trust boundary and routes messages between agents.
//
// The gateway owns no agents. It authenticates them with Noise-KK, decides
// what may reach what according to group membership, and holds messages for
// agents that are not currently connected.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	pub "github.com/waywardgeek/ensemble/agent/gateway"
	gw "github.com/waywardgeek/ensemble/agent/internal/gateway"
)

func main() {
	var (
		configPath = flag.String("config", "gateway.json", "path to the gateway config file")
		keyPath    = flag.String("key", "", "path to the gateway's private key (default <data_dir>/gateway.key)")
		listen     = flag.String("listen", "", "override the listen address from config")
		printKey   = flag.Bool("print-key", false, "print the gateway's public key and exit")
	)
	flag.Parse()

	if err := run(*configPath, *keyPath, *listen, *printKey); err != nil {
		// Exiting lives here, in main, so that everything below is testable
		// and returns errors rather than killing the process.
		fmt.Fprintf(os.Stderr, "gateway: %v\n", err)
		os.Exit(1)
	}
}

func run(configPath, keyPath, listen string, printKey bool) error {
	cfg, err := gw.LoadConfig(configPath)
	if err != nil {
		return err
	}
	if listen != "" {
		cfg.Gateway.Listen = listen
	}
	if keyPath == "" {
		keyPath = filepath.Join(cfg.Gateway.DataDir, "gateway.key")
	}

	key, err := pub.LoadOrCreateKeypair(keyPath)
	if err != nil {
		return err
	}
	if printKey {
		fmt.Println(pub.EncodePublicKey(key.Public))
		return nil
	}

	server, err := gw.NewServer(cfg, key)
	if err != nil {
		return err
	}

	// An operator cannot configure an agent without this, and a gateway
	// whose key is a mystery is a gateway nobody can connect to.
	log.Printf("gateway: public key %s", pub.EncodePublicKey(key.Public))
	log.Printf("gateway: listening on %s", cfg.Gateway.Listen)
	for id := range cfg.Agents {
		log.Printf("gateway: configured agent %s", id)
	}

	httpSrv := &http.Server{
		Addr:              cfg.Gateway.Listen,
		Handler:           server.Handler(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Persist periodically as well as on shutdown. A gateway killed with
	// SIGKILL never runs its shutdown path, and losing an hour of
	// registrations to an abrupt stop is the kind of data loss that is
	// only noticed much later.
	go persistPeriodically(ctx, server)

	err = server.Serve(ctx, httpSrv)
	if err == http.ErrServerClosed {
		err = nil
	}
	if persistErr := server.Persist(); persistErr != nil {
		log.Printf("gateway: final persist failed: %v", persistErr)
	}
	log.Printf("gateway: stopped")
	return err
}

func persistPeriodically(ctx context.Context, server *gw.Server) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := server.Persist(); err != nil {
				log.Printf("gateway: persist failed: %v", err)
			}
		}
	}
}
