package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/keepdevops/cofiswarm-kvpool/internal/bus"
	"github.com/keepdevops/cofiswarm-kvpool/internal/httpapi"
	"github.com/keepdevops/cofiswarm-kvpool/internal/policy"
	"github.com/keepdevops/cofiswarm-observer-sdk/pkg/buspresence"
	"github.com/keepdevops/cofiswarm-observer-sdk/pkg/servicecomponent"
)

func main() {
	addr := flag.String("listen", ":8014", "listen address (HTTP mode)")
	cfgPath := flag.String("config", "", "policy YAML (default: $COFISWARM_KVPOOL_CONFIG or configs/kvpool.yaml)")
	busMode := flag.Bool("bus", false, "serve .kv.* on the NATS observer bus instead of HTTP")
	natsURL := flag.String("nats", "nats://127.0.0.1:4222", "NATS URL (bus mode)")
	flag.Parse()

	path := *cfgPath
	if path == "" {
		if v := os.Getenv("COFISWARM_KVPOOL_CONFIG"); v != "" {
			path = v
		} else {
			path = "configs/kvpool.yaml"
		}
	}
	cfg, err := policy.Load(path)
	if err != nil {
		log.Fatalf("load policy config %s: %v", path, err)
	}

	if *busMode {
		serveBus(*natsURL, cfg)
		return
	}

	// Carrier presence (broker-free, default-off via COFISWARM_BRIDGE_URL): appear in the
	// observer live roster over the zmq-bridge without needing a NATS broker.
	stopPresence := buspresence.StartPresence(os.Getenv("COFISWARM_BRIDGE_URL"), "kvpool", map[string]any{"name": "kvpool"})
	defer stopPresence()

	log.Printf("kvpool listening on %s (enabled=%v proactive=%.2f budgets=%d)",
		*addr, cfg.Enabled, cfg.ProactiveThreshold, len(cfg.Budgets))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	httpSrv := &http.Server{Addr: *addr, Handler: httpapi.New(cfg).Handler()}
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("kvpool http serve: %v", err)
		}
	}()

	<-ctx.Done()
	log.Print("kvpool stopping")
	if err := httpSrv.Shutdown(context.Background()); err != nil {
		log.Printf("kvpool http shutdown: %v", err)
	}
}

func serveBus(url string, cfg policy.Config) {
	nc, err := servicecomponent.Connect(url, "cofiswarm-kvpool")
	if err != nil {
		log.Fatalf("bus connect %s: %v", url, err)
	}
	defer nc.Close()
	comp := servicecomponent.New(nc, "kvpool", "kvpool", bus.Routes(cfg))
	if err := comp.Start(); err != nil {
		log.Fatalf("bus start: %v", err)
	}
	defer comp.Shutdown()
	log.Printf("kvpool on bus %s (.kv.admit/.kv.evaluate/.kv.policy)", url)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Print("kvpool bus stopping")
}
