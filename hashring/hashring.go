// Package hashring provides a consistent hash ring implementation for distributed
// load balancing and caching scenarios.
//
// It uses virtual nodes to improve distribution uniformity and supports
// per-node weights to control traffic allocation.
package hashring

import (
	"crypto/md5"
	"encoding/binary"
	"sort"
	"strconv"
)

const (
	// defaultMultiplier defines how many virtual nodes each unit of weight creates.
	defaultMultiplier = 200
)

// HashRing is a consistent hash ring that maps keys to nodes.
//
// It distributes keys uniformly across nodes using virtual nodes.
// A node with higher weight receives proportionally more virtual nodes.
//
// HashRing is not safe for concurrent use.
type HashRing struct {
	nodes    map[string]int // node name -> weight
	replicas []replica      // sorted virtual nodes
}

// replica represents a single virtual node on the ring.
type replica struct {
	hash uint32
	node string
}

// New creates a HashRing with the given nodes, each assigned a default weight of 1.
//
// Parameters:
//   - nodes: the physical node names to place on the ring.
//
// Returns:
//   - *HashRing: a new hash ring containing the nodes.
func New(nodes ...string) *HashRing {
	weights := make(map[string]int, len(nodes))
	for _, n := range nodes {
		weights[n] = 1
	}
	return NewWeighted(weights)
}

// NewWeighted creates a HashRing with per-node weights.
//
// A node with weight 2 receives twice as many virtual nodes as a node with weight 1,
// and therefore handles roughly twice the traffic.
//
// Parameters:
//   - weights: a map from node name to weight. Weight must be > 0.
//
// Returns:
//   - *HashRing: a new hash ring containing the weighted nodes.
func NewWeighted(weights map[string]int) *HashRing {
	hr := &HashRing{
		nodes: make(map[string]int, len(weights)),
	}
	for node, w := range weights {
		if w <= 0 {
			w = 1
		}
		hr.nodes[node] = w
	}
	hr.rebuild()
	return hr
}

// rebuild regenerates all virtual nodes from the current node weights.
func (hr *HashRing) rebuild() {
	var total int
	for _, w := range hr.nodes {
		total += w * defaultMultiplier
	}

	hr.replicas = make([]replica, 0, total)
	for node, w := range hr.nodes {
		count := w * defaultMultiplier
		for i := 0; i < count; i++ {
			hr.replicas = append(hr.replicas, replica{
				hash: hashKey(node + "#" + strconv.Itoa(i)),
				node: node,
			})
		}
	}

	sort.Slice(hr.replicas, func(i, j int) bool {
		return hr.replicas[i].hash < hr.replicas[j].hash
	})
}

// hashKey computes a 32-bit hash for the given key using MD5.
func hashKey(key string) uint32 {
	sum := md5.Sum([]byte(key))
	return binary.BigEndian.Uint32(sum[:4])
}

// Get returns the node responsible for the given key.
//
// It hashes the key and walks clockwise around the ring to find the first
// virtual node. If the ring is empty, ok is false.
//
// Parameters:
//   - key: the key to look up.
//
// Returns:
//   - string: the selected node name.
//   - bool: true if a node was found.
func (hr *HashRing) Get(key string) (string, bool) {
	if len(hr.replicas) == 0 {
		return "", false
	}
	h := hashKey(key)
	idx := sort.Search(len(hr.replicas), func(i int) bool {
		return hr.replicas[i].hash >= h
	})
	if idx == len(hr.replicas) {
		idx = 0
	}
	return hr.replicas[idx].node, true
}

// GetN returns up to n distinct nodes for the given key, walking clockwise
// around the ring.
//
// If n is greater than the number of physical nodes, all nodes are returned.
// If the ring is empty, an empty slice is returned.
//
// Parameters:
//   - key: the key to look up.
//   - n: the maximum number of distinct nodes to return.
//
// Returns:
//   - []string: the selected node names in ring order.
func (hr *HashRing) GetN(key string, n int) []string {
	if n <= 0 || len(hr.replicas) == 0 {
		return nil
	}
	if n > len(hr.nodes) {
		n = len(hr.nodes)
	}

	h := hashKey(key)
	idx := sort.Search(len(hr.replicas), func(i int) bool {
		return hr.replicas[i].hash >= h
	})

	result := make([]string, 0, n)
	seen := make(map[string]struct{}, n)
	for len(result) < n {
		if idx >= len(hr.replicas) {
			idx = 0
		}
		node := hr.replicas[idx].node
		if _, ok := seen[node]; !ok {
			seen[node] = struct{}{}
			result = append(result, node)
		}
		idx++
	}
	return result
}

// Add inserts a new node with the given weight into the ring.
//
// If the node already exists, its weight is updated and the ring is rebuilt.
//
// Parameters:
//   - node: the node name to add.
//   - weight: the weight to assign. Must be > 0; values <= 0 are treated as 1.
//
// Returns:
//   - *HashRing: the receiver for method chaining.
func (hr *HashRing) Add(node string, weight int) *HashRing {
	if weight <= 0 {
		weight = 1
	}
	hr.nodes[node] = weight
	hr.rebuild()
	return hr
}

// Remove deletes a node from the ring.
//
// Parameters:
//   - node: the node name to remove.
//
// Returns:
//   - bool: true if the node existed and was removed.
func (hr *HashRing) Remove(node string) bool {
	if _, ok := hr.nodes[node]; !ok {
		return false
	}
	delete(hr.nodes, node)
	hr.rebuild()
	return true
}

// Nodes returns a sorted list of all physical nodes currently in the ring.
//
// Returns:
//   - []string: the sorted node names.
func (hr *HashRing) Nodes() []string {
	out := make([]string, 0, len(hr.nodes))
	for n := range hr.nodes {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}
