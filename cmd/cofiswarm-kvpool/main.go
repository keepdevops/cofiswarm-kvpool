package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/keepdevops/cofiswarm-kvpool/internal/httpapi"
)

func main() {
	addr := flag.String("listen", ":8014", "listen address")
	flag.Parse()
	log.Printf("kvpool listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, httpapi.New().Handler()))
}
