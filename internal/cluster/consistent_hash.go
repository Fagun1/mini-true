package cluster

import (
	"fmt"
	"sync"

	"github.com/spaolacci/murmur3"
)

type ConsistentHashRing struct {
	ring         map[uint32]string
	sortedHashes []uint32
	virtualNodes int
	nodes        map[string]bool
	hashFunc     func([]byte) uint32
	mu           sync.RWMutex
}

func NewConsistentHashRing(virtualNodes int) *ConsistentHashRing {
	if virtualNodes <= 0 {
		virtualNodes = 150
	}

	return &ConsistentHashRing{
		ring:         make(map[uint32]string),
		sortedHashes: make([]uint32, 0),
		virtualNodes: virtualNodes,
		nodes:        make(map[string]bool),
		// Changing the hash function changes key ownership and requires a full
		// cluster data rebalance before a live upgrade. This is expected to be
		// paired with the segment migration path before production rollout.
		hashFunc: murmur3.Sum32,
	}
}

func (chr *ConsistentHashRing) hashKey(key string) uint32 {
	return chr.hashFunc([]byte(key))
}


