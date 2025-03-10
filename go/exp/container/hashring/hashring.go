// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hashring provides a consistent hashing function.
//
// HashRing, aka NodeLocator hashing, is often used to distribute requests to a changing set of servers.
// For example, say you have some cache servers cacheA, cacheB, and cacheC.
// You want to decide which cache server to use to look up information on a user.
//
// You could use a typical hash table and hash the user id
// to one of cacheA, cacheB, or cacheC.  But with a typical hash table, if you add or remove a server,
// almost all keys will get remapped to different results, which basically could bring your service
// to a grinding halt while the caches get rebuilt.
//
// With a consistent hash, adding or removing a server drastically reduces the number of keys that
// get remapped.
//
// Read more about consistent hashing on wikipedia:  http://en.wikipedia.org/wiki/Consistent_hashing
package hashring

import (
	"iter"
	"maps"
	"math"
	"slices"
)

const defaultNumReps = 160

// HashRing holds the information about the allNodes of the consistent hash nodes.
//
// Node represents a node in the consistent hash ring.
// {}	-> 127.0.0.1:11311 -> 127.0.0.1:11311-0 -> 1234
// Node ->       Key       ->     IterateKey    -> HashKey
//
//go:generate go-option -type "HashRing"
type HashRing[Node comparable] struct {
	// The List of nodes to use in the Ketama consistent hash continuum
	//
	// This simulates the structure of keys used in the Ketama consistent hash ring,
	// which stores the virtual node HashKeys on the physical nodes.
	// All nodes in the cluster are topped by virtual nodes.
	// In principle, it is a brute-force search to find the first complete HashKey
	//
	// For example,
	// Node ->       Key       ->      IterateKey     -> HashKey
	// {}	-> 127.0.0.1:11311 -> 127.0.0.1:11311-0   ->  1234
	// {}	-> 127.0.0.1:11311 -> 127.0.0.1:11311-160 ->  256
	// {}	-> 127.0.0.1:11311 -> 127.0.0.1:11311-320 ->  692
	sortedKeys []uint32          // []HashKey, Index for nodes binary search
	nodeByKey  map[uint32]Node   // <HashKey,Node>
	allNodes   map[Node]struct{} // <Node>

	// The hash algorithm to use when choosing a node in the Ketama consistent hash continuum
	hashAlg HashAlgorithm

	// node weights for ketama, a map from InetSocketAddress to weight as Integer
	weightByNode map[Node]int
	isWeighted   bool

	// the number of discrete hashes that should be defined for each node in the continuum.
	numReps int
	// the format used to name the nodes in Ketama, either SpyMemcached or LibMemcached
	nodeKeyFormatter Formatter[Node]
}

// New creates a hash ring of n replicas for each entry.
func New[Node comparable](opts ...HashRingOption[Node]) *HashRing[Node] {
	r := &HashRing[Node]{
		nodeByKey:        make(map[uint32]Node),
		allNodes:         make(map[Node]struct{}),
		hashAlg:          KetamaHash,
		weightByNode:     make(map[Node]int),
		numReps:          defaultNumReps,
		nodeKeyFormatter: NewKetamaNodeKeyFormatter[Node](SpyMemcached),
	}
	r.ApplyOptions(opts...)
	if r.isWeighted && len(r.weightByNode) == 0 {
		r.isWeighted = false
	}

	return r
}

// AddNodes inserts nodes into the consistent hash cycle.
func (c *HashRing[Node]) AddNodes(nodes ...Node) {
	if c.isWeighted {
		c.addWeightNodes(nodes...)
		return
	}
	c.addNoWeightNodes(nodes...)
}

// SetNodes setups the HashRing with the list of nodes it should use.
// If there are existing nodes not present in nodes, they will be removed.
// @param nodes a List of Nodes for this HashRing to use in
// its continuum
func (c *HashRing[Node]) SetNodes(nodes ...Node) {
	if c.isWeighted {
		c.setWeightNodes(nodes...)
		return
	}
	c.setNoWeightNodes(nodes...)
}

// RemoveAllNodes removes all nodes in the continuum.
func (c *HashRing[Node]) RemoveAllNodes() {
	c.sortedKeys = nil
	c.nodeByKey = make(map[uint32]Node)
	c.allNodes = make(map[Node]struct{})
}

// Get returns an element close to where name hashes to in the nodes.
func (c *HashRing[Node]) Get(name string) (Node, bool) {
	if len(c.nodeByKey) == 0 {
		var zeroN Node
		return zeroN, false
	}
	return c.getPrimaryNode(name)
}

// GetSince returns an iterator over distinct nodes in hashring, start from where name hashes to in the nodes.
func (c *HashRing[Node]) GetSince(name string) iter.Seq[Node] {
	return func(yield func(Node) bool) {
		if len(c.nodeByKey) == 0 {
			return
		}

		firstKey, found := c.tailSearch(c.getHashKey(name))
		if !found {
			firstKey = 0
		}
		firstNode := c.getNodeByHashKeyIndex(firstKey)

		if !yield(firstNode) {
			return
		}

		nodesSet := make(map[Node]struct{}) // for unique only
		nodesSet[firstNode] = struct{}{}

		start := firstKey
		var secondNode Node
		for i := start + 1; i != start; i++ {
			if i >= len(c.sortedKeys) {
				i = 0
				i-- // take care of i++ after this loop of for
				continue
			}
			secondNode = c.getNodeByHashKeyIndex(i)
			if _, has := nodesSet[secondNode]; !has {
				if !yield(secondNode) {
					return
				}
				nodesSet[secondNode] = struct{}{}
			}
		}
		return
	}
}

// All returns an iterator over all nodes in hashring.
// If c is empty, the sequence is empty: there is no empty element in the sequence.
func (c *HashRing[Node]) All() iter.Seq[Node] {
	return maps.Keys(c.allNodes)
}

// getAllNodes returns all available nodes
func (c *HashRing[Node]) getAllNodes() []Node {
	return slices.Collect(maps.Keys(c.allNodes))
}

// getPrimaryNode returns the first available node for a name, such as “127.0.0.1:11311-0” for "Alice"
func (c *HashRing[Node]) getPrimaryNode(name string) (Node, bool) {
	return c.getNodeByHashKey(c.getHashKey(name))
}

// getMaxHashKey returns the last available node's HashKey
// that is, Maximum HashKey in the Hash Cycle
func (c *HashRing[Node]) getMaxHashKey() (key uint32, ok bool) {
	if len(c.sortedKeys) == 0 {
		return 0, false
	}
	return c.sortedKeys[len(c.sortedKeys)-1], true
}

// getNodeByHashKey returns the first available node since iterateHashKey, such as HASH(“127.0.0.1:11311-0”)
func (c *HashRing[Node]) getNodeByHashKey(hash uint32) (Node, bool) {
	if len(c.sortedKeys) == 0 {
		var zeroN Node
		return zeroN, false
	}

	rv, has := c.nodeByKey[hash]
	if has {
		return rv, true
	}
	firstKey, found := c.tailSearch(hash)
	if !found {
		firstKey = 0
	}
	return c.getNodeByHashKeyIndex(firstKey), true
}

// getNodeByHashKeyIndex returns the node by index of sorted hash keys.
func (c *HashRing[Node]) getNodeByHashKeyIndex(keyIndex int) (node Node) {
	return c.nodeByKey[c.sortedKeys[keyIndex]]
}

// updateLocator reconstructs the hash ring with the input nodes
func (c *HashRing[Node]) updateLocator(nodes ...Node) {
	c.SetNodes(nodes...)
}

// getNodeRepetitions returns the number of discrete hashes that should be defined for each node
// in the continuum.
func (c *HashRing[Node]) getNodeRepetitions() int {
	return c.numReps
}

// setNoWeightNodes sets all the elements in the hash.
func (c *HashRing[Node]) setNoWeightNodes(nodes ...Node) {
	// Set sets all the elements in the hash.
	// If there are existing elements not present in nodes, they will be removed.
	var nodesToBeRemoved []Node
	// remove missing Nodes
	for k := range c.allNodes {
		var found bool
		for _, v := range nodes {
			if c.isSameNode(k, v) {
				// found
				found = true
				break
			}
		}
		if !found {
			nodesToBeRemoved = append(nodesToBeRemoved, k)
		}
	}
	if len(nodesToBeRemoved) == len(nodes) {
		c.RemoveAllNodes()
	} else {
		c.removeNoWeightNodes(nodesToBeRemoved...)
	}
	// add all missing elements present in nodes.
	var nodesToBeAdded []Node
	for _, k := range nodes {
		var found bool
		for v := range c.allNodes {
			if c.isSameNode(k, v) {
				found = true
				break
			}
		}
		if !found {
			nodesToBeAdded = append(nodesToBeAdded, k)
		}
	}
	c.addNoWeightNodes(nodesToBeAdded...)
}

// setWeightNodes sets all the elements in the hash.
func (c *HashRing[Node]) setWeightNodes(nodes ...Node) {
	c.RemoveAllNodes()
	numReps := c.getNodeRepetitions()
	nodeCount := len(nodes)
	totalWeight := 0

	for _, node := range nodes {
		totalWeight += c.weightByNode[node]
	}

	// add all elements present in nodes.
	for _, node := range nodes {
		thisWeight := c.weightByNode[node]
		percent := float64(thisWeight) / float64(totalWeight)
		// floor(percent * numReps * nodeCount + 1e10)
		pointerPerServer := (int)(math.Floor(percent*(float64(numReps))*float64(nodeCount) + 1e10))
		c.addNodeWithoutSort(node, pointerPerServer)
	}

	// sort keys
	c.updateSortedNodes()
}

// addWeightNodes adds a node to the hash without sorting the keys.
func (c *HashRing[Node]) addWeightNodes(nodes ...Node) {
	c.setWeightNodes(append(c.getAllNodes(), nodes...)...)
}

// addNoWeightNodes adds a node to the hash without sorting the keys.
func (c *HashRing[Node]) addNoWeightNodes(nodes ...Node) {
	numReps := c.getNodeRepetitions()

	for _, node := range nodes {
		c.addNodeWithoutSort(node, numReps)
	}

	c.updateSortedNodes()
}

// addNodeWithoutSort adds a node to the hash without sorting the keys.
func (c *HashRing[Node]) addNodeWithoutSort(node Node, numReps int) {
	// Ketama does some special work with md5 where it reuses chunks.
	// Check to be backwards compatible, the hash algorithm does not
	// matter for Ketama, just the placement should always be done using
	// MD5

	// KETAMA_HASH, Special Case, batch mode to speedup

	for i := 0; i < numReps; {
		positions := c.getIterateHashKeyForNode(node, i)
		if len(positions) == 0 {
			numReps++
			i++ // ignore no hash node
			break
		}

		for j, pos := range positions {
			if i+j > numReps { // out of bound
				break
			}
			if _, has := c.nodeByKey[pos]; has {
				// skip this node, duplicated
				numReps++
				continue
			}
			c.nodeByKey[pos] = node
		}
		i += len(positions)
	}

	c.allNodes[node] = struct{}{}
}

// RemoveNodes removes nodes from the consistent hash cycle
func (c *HashRing[Node]) RemoveNodes(nodes ...Node) {
	if c.isWeighted {
		c.removeWeightNodes(nodes...)
		return
	}
	c.removeNoWeightNodes(nodes...)
}

// removeWeightNodes removes nodes from the consistent hash cycle
func (c *HashRing[Node]) removeWeightNodes(nodes ...Node) {
	for _, node := range nodes {
		delete(c.allNodes, node)
	}
	c.setWeightNodes(c.getAllNodes()...)
}

func (c *HashRing[Node]) removeNoWeightNodes(nodes ...Node) {
	numReps := c.getNodeRepetitions()

	for _, node := range nodes {
		for i := 0; i < numReps; {
			positions := c.getIterateHashKeyForNode(node, i)
			if len(positions) == 0 {
				// ignore no hash node
				numReps++
				i++
				continue
			}

			for j, pos := range positions {
				if i+j > numReps { // out of bound
					break
				}
				if n, has := c.nodeByKey[pos]; has {
					if !c.isSameNode(n, node) {
						numReps++ // ignore no hash node
						continue
					}
					delete(c.nodeByKey, pos)
				}
			}
			i += len(positions)
		}
		delete(c.allNodes, node)
	}
	c.updateSortedNodes()
}

// tailSearch returns the first available node since iterateHashKey's Index, such as Index(HASH(“127.0.0.1:11311-0”))
func (c *HashRing[Node]) tailSearch(key uint32) (i int, found bool) {
	// Search uses binary search to find and return the smallest index since iterateHashKey's Index
	return slices.BinarySearchFunc(c.sortedKeys, key, func(v uint32, key uint32) int {
		if v >= key {
			return 0
		}
		return -1
	})
}

// updateSortedNodes sorts the keys in ascending order.
func (c *HashRing[Node]) updateSortedNodes() {
	hashes := c.sortedKeys[:0]
	// reallocate if we're holding on to too much (1/4th)
	// len(nodes) * replicas < cap / 4
	// len(c.nodeByKey) ≈ len(c.allNodes)*c.numReps
	if cap(c.sortedKeys)/4 > len(c.nodeByKey) {
		hashes = nil
	}
	for k := range c.nodeByKey {
		hashes = append(hashes, k)
	}
	slices.Sort(hashes)
	c.sortedKeys = hashes
}

// isSameNode checks if two nodes are the same by the key.
func (c *HashRing[Node]) isSameNode(n1, n2 Node) bool {
	return c.nodeKeyFormatter.FormatNodeKey(n1, 0) == c.nodeKeyFormatter.FormatNodeKey(n2, 0)
}
