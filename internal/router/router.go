// Package router accepts sensor data and resolves its owning hash-ring node.
package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/minitrue/internal/cluster"
	"github.com/minitrue/internal/models"
)

// Router is the HTTP ingestion gateway for the early MiniTrue pipeline.
type Router struct {
	ring *cluster.ConsistentHashRing
}

func New(ring *cluster.ConsistentHashRing) *Router {
	return &Router{ring: ring}
}

// ServeHTTP accepts POST /route and reports the node selected for the datapoint.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var point models.DataPoint
	body, err := io.ReadAll(req.Body)
	if err != nil || json.Unmarshal(body, &point) != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if point.DeviceID == "" || point.MetricName == "" {
		http.Error(w, "missing device_id or metric_name", http.StatusBadRequest)
		return
	}
	if point.Timestamp == 0 {
		point.Timestamp = time.Now().Unix()
	}

	node, err := r.Route(point)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "accepted",
		"node":       node,
		"data_point": point,
	})
}

// Route builds the composite device/metric key and resolves its owner.
func (r *Router) Route(point models.DataPoint) (string, error) {
	if r.ring == nil {
		return "", fmt.Errorf("hash ring is not initialised")
	}
	key := point.DeviceID + ":" + point.MetricName
	nodes, err := r.ring.GetNodes(key, 1)
	if err != nil || len(nodes) == 0 {
		return "", fmt.Errorf("no node available for key %q", key)
	}
	return nodes[0], nil
}
