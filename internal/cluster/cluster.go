package cluster

var hashRing = NewConsistentHashRing(150)

func GetNodesForKey(key string, replicationFactor int) []string {
	nodes, err := hashRing.GetNodes(key, replicationFactor)
	if err != nil || len(nodes) == 0 {
		return nil
	}
	return nodes
}

func GetHashRing() *ConsistentHashRing {
	return hashRing
}

func SetHashRing(ring *ConsistentHashRing) {
	hashRing = ring
}


