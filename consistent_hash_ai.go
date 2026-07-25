package main

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strconv"
	"sync"
)

type HashFunc func(data []byte) uint32

type ConsistentHash struct {
	hash     HashFunc
	replicas int
	ring     []uint32          // Sorted virtual node hashes
	nodeMap  map[uint32]string // Hash -> Physical Node Name
	storage  map[string][]string // Physical Node Name -> Data Items
	mu       sync.RWMutex
}

func NewConsistentHash(replicas int, fn HashFunc) *ConsistentHash {
	if fn == nil {
		fn = func(data []byte) uint32 {
			h := fnv.New32a()
			h.Write(data)
			return h.Sum32()
		}
	}
	return &ConsistentHash{
		replicas: replicas,
		hash:     fn,
		nodeMap:  make(map[uint32]string),
		storage:  make(map[string][]string),
	}
}

// Add inserts physical nodes into the ring and initializes their storage.
func (c *ConsistentHash) Add(nodes ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, node := range nodes {
		if _, exists := c.storage[node]; !exists {
			c.storage[node] = make([]string, 0)
		}

		for i := 0; i < c.replicas; i++ {
			vNodeKey := node + "#" + strconv.Itoa(i)
			hash := c.hash([]byte(vNodeKey))
			if _, exists := c.nodeMap[hash]; !exists {
				c.ring = append(c.ring, hash)
			}
			c.nodeMap[hash] = node
		}
	}

	sort.Slice(c.ring, func(i, j int) bool {
		return c.ring[i] < c.ring[j]
	})
}

// GetNodeForKey returns the physical node responsible for a given key (unlocked internal helper).
func (c *ConsistentHash) getNodeForKey(key string) string {
	if len(c.ring) == 0 {
		return ""
	}
	hash := c.hash([]byte(key))
	idx := sort.Search(len(c.ring), func(i int) bool {
		return c.ring[i] >= hash
	})
	if idx == len(c.ring) {
		idx = 0
	}
	return c.nodeMap[c.ring[idx]]
}

// AddData hashes the data item, balances it to the correct node, and stores it.
func (c *ConsistentHash) AddData(data string) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	targetNode := c.getNodeForKey(data)
	if targetNode == "" {
		return ""
	}

	c.storage[targetNode] = append(c.storage[targetNode], data)
	return targetNode
}

// Remove deletes a physical node and transfers all its stored data 
// to the next adjacent physical node on the ring for each item.
func (c *ConsistentHash) Remove(node string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 1. Extract data stored on the node being removed
	migratingData := c.storage[node]
	delete(c.storage, node)

	// 2. Remove all virtual node replicas from the ring
	for i := 0; i < c.replicas; i++ {
		vNodeKey := node + "#" + strconv.Itoa(i)
		hash := c.hash([]byte(vNodeKey))
		delete(c.nodeMap, hash)

		idx := sort.Search(len(c.ring), func(j int) bool {
			return c.ring[j] >= hash
		})
		if idx < len(c.ring) && c.ring[idx] == hash {
			c.ring = append(c.ring[:idx], c.ring[idx+1:]...)
		}
	}

	// 3. Re-route data to their new successor nodes on the updated ring
	if len(c.ring) > 0 {
		for _, item := range migratingData {
			newNode := c.getNodeForKey(item)
			c.storage[newNode] = append(c.storage[newNode], item)
		}
	}
}

// GetStorage returns a snapshot of data distribution across all physical nodes.
func (c *ConsistentHash) GetStorage() map[string][]string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshot := make(map[string][]string)
	for k, v := range c.storage {
		copied := make([]string, len(v))
		copy(copied, v)
		snapshot[k] = copied
	}
	return snapshot
}


func main() {
	ring := NewConsistentHash(3, nil)
	ring.Add("node-A", "node-B", "node-C")

	// 1. Load balance and insert data items
	items := []string{"user_1", "user_2", "order_99", "item_404", "session_7", "token_x"}
	fmt.Println("--- Storing Data Items ---")
	for _, item := range items {
		assignedNode := ring.AddData(item)
		fmt.Printf("Item '%s' stored in -> %s\n", item, assignedNode)
	}

	fmt.Println("\n--- Node Storage Before Removal ---")
	printStorage(ring.GetStorage())

	// 2. Remove node-B (its data will re-route to remaining nodes)
	fmt.Println("\n--- Removing 'node-B' ---")
	ring.Remove("node-B")

	fmt.Println("\n--- Node Storage After Removal & Migration ---")
	printStorage(ring.GetStorage())
}

func printStorage(storage map[string][]string) {
	for node, data := range storage {
		fmt.Printf("%s: %v\n", node, data)
	}
}
