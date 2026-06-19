package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/keepdevops/cofiswarm-kvpool/internal/httpapi"
	"github.com/keepdevops/cofiswarm-kvpool/internal/policy"
)

func main() {
	addr := flag.String("listen", ":8014", "listen address")
	cfgPath := flag.String("config", "", "policy YAML (default: $COFISWARM_KVPOOL_CONFIG or configs/kvpool.yaml)")
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
	log.Printf("kvpool listening on %s (enabled=%v proactive=%.2f budgets=%d)",
		*addr, cfg.Enabled, cfg.ProactiveThreshold, len(cfg.Budgets))
	log.Fatal(http.ListenAndServe(*addr, httpapi.New(cfg).Handler()))
}
