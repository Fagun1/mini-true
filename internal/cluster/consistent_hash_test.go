package cluster

import "testing"

func TestGetNodesEmptyRing(t *testing.T) {
	ring := NewConsistentHashRing(10)

	if _, err := ring.GetNodes("sensor_1:temperature", 1); err == nil {
		t.Fatal("expected an empty ring lookup to fail")
	}
}

func TestGetNodesSingleNode(t *testing.T) {
	ring := NewConsistentHashRing(10)
	ring.AddNode("polaris")

	nodes, err := ring.GetNodes("sensor_1:temperature", 1)
	if err != nil {
		t.Fatalf("single-node lookup failed: %v", err)
	}
	if len(nodes) != 1 || nodes[0] != "polaris" {
		t.Fatalf("expected polaris, got %v", nodes)
	}
}

func TestGetNodesReturnsDistinctNodes(t *testing.T) {
	ring := NewConsistentHashRing(20)
	for _, nodeID := range []string{"polaris", "sirius", "vega"} {
		ring.AddNode(nodeID)
	}

	nodes, err := ring.GetNodes("sensor_1:temperature", 3)
	if err != nil {
		t.Fatalf("multi-node lookup failed: %v", err)
	}
	if len(nodes) != 3 {
		t.Fatalf("expected three nodes, got %v", nodes)
	}

	seen := make(map[string]bool)
	for _, nodeID := range nodes {
		if seen[nodeID] {
			t.Fatalf("node %q returned more than once: %v", nodeID, nodes)
		}
		seen[nodeID] = true
	}
}

func TestRemoveNodeUpdatesRing(t *testing.T) {
	ring := NewConsistentHashRing(10)
	ring.AddNode("polaris")
	ring.AddNode("sirius")
	ring.RemoveNode("sirius")

	nodes, err := ring.GetNodes("sensor_1:temperature", 2)
	if err != nil {
		t.Fatalf("lookup after removal failed: %v", err)
	}
	if len(nodes) != 1 || nodes[0] != "polaris" {
		t.Fatalf("expected only polaris after removal, got %v", nodes)
	}
}
