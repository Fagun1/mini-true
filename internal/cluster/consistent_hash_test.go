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

func TestAddNodeDoesNotDuplicateNode(t *testing.T) {
	ring := NewConsistentHashRing(10)
	ring.AddNode("polaris")
	ring.AddNode("polaris")

	if len(ring.nodes) != 1 {
		t.Fatalf("expected one physical node, got %d", len(ring.nodes))
	}
	if len(ring.sortedHashes) != ring.virtualNodes {
		t.Fatalf("expected %d virtual hashes, got %d", ring.virtualNodes, len(ring.sortedHashes))
	}
}

func TestGetNodesWrapsToFirstHash(t *testing.T) {
	ring := NewConsistentHashRing(10)
	ring.AddNode("polaris")
	ring.AddNode("sirius")

	firstHash := ring.sortedHashes[0]
	var key string
	for i := 0; i < 1000000; i++ {
		candidate := "wrap-key-" + string(rune(i))
		if ring.hashKey(candidate) > ring.sortedHashes[len(ring.sortedHashes)-1] {
			key = candidate
			break
		}
	}
	if key == "" {
		t.Fatal("could not find a key beyond the last ring position")
	}

	nodes, err := ring.GetNodes(key, 1)
	if err != nil {
		t.Fatalf("wrap-around lookup failed: %v", err)
	}
	if len(nodes) != 1 || ring.ring[firstHash] != nodes[0] {
		t.Fatalf("expected lookup to wrap to %q, got %v", ring.ring[firstHash], nodes)
	}
}
