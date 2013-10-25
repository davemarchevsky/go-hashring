package hashring

import (
	"fmt"
	"math"
	"testing"
)

func _getRandomStrings(n int) []string {
	words := make([]string, n)
	for i := 0; i < n; i++ {
		words[i] = fmt.Sprintf("key %d of %d", i, n)
	}
	return words
}

func TestDistributionWithoutWeights(t *testing.T) {
	// N random strings, M nodes:
	// ensure that they are relatively evenly distributed among nodes
	nodes := []string{}
	for i := 0; i < 12; i++ {
		nodes = append(nodes, fmt.Sprintf("host%2d", i))
	}
	ring := New(nodes)
	num_keys := 50000
	keys := _getRandomStrings(num_keys)
	_checkNodeBalance(t, ring, keys)
}

func _checkNodeBalance(t *testing.T, ring *HashRing, keys []string) {
	nodes := ring.Nodes
	num_keys := len(keys)
	ideal_count := int(float32(num_keys) / float32(len(nodes)))
	expected_variation := 0.25 // Usually more like 10%
	delta := int(expected_variation * float64(ideal_count))

	countByNode := make(map[string]int)
	for _, key := range keys {
		node := ring.GetNode(key)
		countByNode[node] += 1
	}

	for node, count := range countByNode {
		//fmt.Printf("Node %s: %d\n", node, count)
		if int(math.Abs(float64(count-ideal_count))) > delta {
			t.Errorf("%s not well balanced: got %d, expected %d +/- %d", node, count, ideal_count, delta)
		}
	}
}

func TestDistributionAfterAddingNode(t *testing.T) {
	// 	* N random strings, M nodes, then add a node:
	// ensure that the number of keys that switch nodes is sane.
	nodes := []string{"host1", "host2", "host3", "host4", "host5"}
	ring := New(nodes)
	originalNodesByKey := make(map[string]string)
	num_keys := 50000
	keys := _getRandomStrings(num_keys)
	for _, key := range keys {
		node := ring.GetNode(key)
		originalNodesByKey[key] = node
	}

	// Add a node.
	nodes = append(nodes, "host6")
	ring = New(nodes)

	// Check how many keys moved nodes.
	moved_keys := 0
	expected_to_move := int(float32(num_keys) / float32(len(nodes)))
	expected_variation := 0.25
	delta := int(expected_variation * float64(expected_to_move))
	for _, key := range keys {
		node := ring.GetNode(key)
		if node != originalNodesByKey[key] {
			moved_keys += 1
		}
	}
	if int(math.Abs(float64(expected_to_move-moved_keys))) > delta {
		t.Errorf(
			"%d out of %d moved, expected %d +/- %d",
			moved_keys, num_keys, expected_to_move, delta)
	}
	// Distribution should still be good.
	_checkNodeBalance(t, ring, keys)

}

func TestDistributionAfterRemovingNode(t *testing.T) {
	// * N random strings, M nodes, then remove a node:
	//ensure that the number of keys that move nodes is sane
	nodes := []string{
		"host1", "host2", "host3", "host4", "host5",
	}
	ring := New(nodes)
	originalNodesByKey := make(map[string]string)
	countByNode := make(map[string]int)
	num_keys := 10000
	keys := _getRandomStrings(num_keys)
	for _, key := range keys {
		node := ring.GetNode(key)
		originalNodesByKey[key] = node
		countByNode[node] += 1
	}

	// Remove a node.
	expected_to_move := countByNode["host1"]
	nodes = nodes[1:]
	ring = New(nodes)

	// Check how many keys moved nodes. Should be exact.
	moved_keys := 0
	for _, key := range keys {
		node := ring.GetNode(key)
		if node != originalNodesByKey[key] {
			moved_keys += 1
		}
	}
	if expected_to_move != moved_keys {
		t.Errorf(
			"%d out of %d moved, expected %d",
			moved_keys, num_keys, expected_to_move)
	}
	// Distribution should still be good.
	_checkNodeBalance(t, ring, keys)

}

func TestDistributionWithWeights(t *testing.T) {
	// N random strings, M nodes:
	// ensure that they are relatively evenly distributed among nodes
	nodeMap := map[string]Weight{
		"little": 1,
		"medium": 2,
		"big":    3,
	}
	ring := NewWithWeights(nodeMap)
	num_keys := 60000
	keys := _getRandomStrings(num_keys)
	countByNode := make(map[string]int)
	for _, key := range keys {
		node := ring.GetNode(key)
		countByNode[node] += 1
	}
	// TODO verify that they're very roughly proportional to the weights.
	if !(countByNode["little"] < countByNode["medium"] && countByNode["medium"] < countByNode["big"]) {
		t.Errorf(
			"Weighting failed: expected %d < %d < %d",
			countByNode["little"],
			countByNode["medium"],
			countByNode["big"],
		)
	}

}
