/*
 This package implements consistent hashing in Go, using md5 as a hashing function.
 It's based off of a python hash_ring implementation (https://pypi.python.org/pypi/hash_ring/1.3.1).
 The functionality of this hashring should be identical to the hash_ring above. 
*/
package hashring

import (
	"crypto/md5"
	"fmt"
	"io"
	"sort"
)

type Weight uint
type RingKey uint32 // hashValGenerate outputs ints > int32 range
type RingKeys []RingKey

type HashRing struct {
	Ring       map[RingKey]string
	Nodes      map[string]Weight
	SortedKeys RingKeys
}

// Satisfy sort.Interface
func (r RingKeys) Len() int {
	return len(r)
}

func (r RingKeys) Less(i, j int) bool {
	return r[i] < r[j]
}

func (r RingKeys) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// Convenience constructor for HashRings with all nodes weighted equally
func New(nodes []string) *HashRing {
	nodeMap := make(map[string]Weight)
	for i := range nodes {
		nodeMap[nodes[i]] = 1
	}
	return NewWithWeights(nodeMap)
}

func NewWithWeights(nodes map[string]Weight) *HashRing {
	r := new(HashRing)
	r.Nodes = nodes
	r.Ring = make(map[RingKey]string)

	r.generateCircle()
	return r
}

func (r *HashRing) generateCircle() {
	var totalWeight Weight = 0
	for _, weight := range r.Nodes {
		totalWeight += weight
	}

	for node, weight := range r.Nodes {
		factor := (40 * len(r.Nodes) * int(weight)) / int(totalWeight)

		for j := 0; j < factor; j++ {
			bKey := hashDigest(fmt.Sprintf("%s-%d", node, j))
			for i := 0; i < 3; i++ {
				offset := i * 4
				key := hashValGenKey(bKey[offset : offset+4])
				r.Ring[key] = node
				r.SortedKeys = append(r.SortedKeys, key)
			}
		}
	}
	sort.Sort(r.SortedKeys)
}

// Given a key, return which node the key hashes to 
func (r *HashRing) GetNode(key string) string {
	pos := r.getNodePos(key)
	return r.Ring[r.SortedKeys[pos]]
}

func (r *HashRing) getNodePos(key string) int {
	nodeKey := genKey(key)
	pos := bisectRightDefault(r.SortedKeys, nodeKey)

	if pos == len(r.SortedKeys) {
		return 0
	}
	return pos
}

/* Reimplementation of the bisect_right method from http://docs.python.org/2/library/bisect.html
   Given a sorted list of RingKeys and a RingKey x, return the position in the sorted list x would have
   if it were inserted.
*/
func bisectRight(list []RingKey, x RingKey, lo, hi int) int {
	for {
		if !(lo < hi) {
			break
		}
		mid := (lo + hi) / 2
		if x < list[mid] {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

/* Most of the time we're doing bisectRight we want:
   lo = beginning of slice
   hi = end of slice
*/
func bisectRightDefault(list []RingKey, x RingKey) int {
	return bisectRight(list, x, 0, len(list))
}

func genKey(keyString string) RingKey {
	bKey := hashDigest(keyString)
	return hashValGenKey(bKey)
}

func hashValGenKey(bKey []byte) RingKey {
	return ((RingKey(bKey[3]) << 24) |
		(RingKey(bKey[2]) << 16) |
		(RingKey(bKey[1]) << 8) |
		(RingKey(bKey[0]) << 0))
}

func hashDigest(key string) []byte {
	m := md5.New()
	io.WriteString(m, key)
	return m.Sum(nil)
}
