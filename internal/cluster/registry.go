package cluster

import (
	"errors"
	"sort"
	"sync"
	"time"
)

var (
	ErrNodeAlreadyExists = errors.New("node already exists")
	ErrNodeNotFound      = errors.New("node not found")
)

type NodeRegistry struct {
	mu    sync.RWMutex
	nodes map[string]Node
}

func NewNodeRegistry() *NodeRegistry {
	return &NodeRegistry{
		nodes: make(map[string]Node),
	}
}

func (r *NodeRegistry) Register(node Node) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.nodes[node.ID]; exists {
		return ErrNodeAlreadyExists
	}

	if r.nodes == nil {
		r.nodes = make(map[string]Node)
	}

	r.nodes[node.ID] = node
	return nil
}

func (r *NodeRegistry) Get(nodeID string) (Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	node, exists := r.nodes[nodeID]
	if !exists {
		return Node{}, ErrNodeNotFound
	}

	return node, nil
}

func (r *NodeRegistry) List() []Node {
	r.mu.RLock()

	nodes := make([]Node, 0, len(r.nodes))
	for _, node := range r.nodes {
		nodes = append(nodes, node)
	}

	r.mu.RUnlock()

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].ID < nodes[j].ID
	})

	return nodes
}

func (r *NodeRegistry) Exists(nodeID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.nodes[nodeID]
	return exists
}

func (r *NodeRegistry) UpdateHeartbeat(nodeID string, heartbeatAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	node, exists := r.nodes[nodeID]
	if !exists {
		return ErrNodeNotFound
	}

	node.LastHeartbeat = heartbeatAt
	r.nodes[nodeID] = node

	return nil
}
