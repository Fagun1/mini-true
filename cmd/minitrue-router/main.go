// minitrue-router runs the early MiniTrue HTTP ingestion and routing service.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/minitrue/internal/cluster"
	"github.com/minitrue/internal/router"
)

func main() {
	port := flag.Int("port", envInt("MINITRUE_ROUTER_PORT", 7070), "HTTP port")
	nodes := flag.String("nodes", envStr("MINITRUE_NODES", "polaris,sirius,vega"), "comma-separated ring nodes")
	flag.Parse()

	ring := cluster.NewConsistentHashRing(150)
	for _, nodeID := range strings.Split(*nodes, ",") {
		if nodeID = strings.TrimSpace(nodeID); nodeID != "" {
			ring.AddNode(nodeID)
		}
	}

	mux := http.NewServeMux()
	mux.Handle("/route", router.New(ring))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok", "time": time.Now().Unix()})
	})

	address := fmt.Sprintf(":%d", *port)
	log.Printf("MiniTrue router listening on %s", address)
	log.Fatal(http.ListenAndServe(address, mux))
}

func envStr(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	var value int
	if _, err := fmt.Sscanf(envStr(key, ""), "%d", &value); err == nil {
		return value
	}
	return fallback
}
