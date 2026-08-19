package cluster

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func testNode(id string) Node {
	now := time.Date(2026, time.August, 19, 12, 0, 0, 0, time.UTC)

	return Node{
		ID:            id,
		Hostname:      id + "-host",
		CPUCores:      8,
		MemoryMB:      16384,
		Status:        NodeStatusHealthy,
		RegisteredAt:  now,
		LastHeartbeat: now,
	}
}

func TestNodeRegistryRegisterAndGet(t *testing.T) {
	registry := NewNodeRegistry()
	node := testNode("node-01")

	if err := registry.Register(node); err != nil {
		t.Fatalf("register node: %v", err)
	}

	got, err := registry.Get(node.ID)
	if err != nil {
		t.Fatalf("get node: %v", err)
	}

	if got != node {
		t.Fatalf("expected %+v, got %+v", node, got)
	}
}

func TestNodeRegistryRejectsDuplicateNode(t *testing.T) {
	registry := NewNodeRegistry()
	node := testNode("node-01")

	if err := registry.Register(node); err != nil {
		t.Fatalf("register node: %v", err)
	}

	err := registry.Register(node)

	if err != ErrNodeAlreadyExists {
		t.Fatalf("expected ErrNodeAlreadyExists, got %v", err)
	}
}

func TestNodeRegistryGetUnknownNode(t *testing.T) {
	registry := NewNodeRegistry()

	_, err := registry.Get("missing-node")

	if err != ErrNodeNotFound {
		t.Fatalf("expected ErrNodeNotFound, got %v", err)
	}
}

func TestNodeRegistryExists(t *testing.T) {
	registry := NewNodeRegistry()
	node := testNode("node-01")

	if registry.Exists(node.ID) {
		t.Fatal("node should not exist before registration")
	}

	if err := registry.Register(node); err != nil {
		t.Fatalf("register node: %v", err)
	}

	if !registry.Exists(node.ID) {
		t.Fatal("node should exist after registration")
	}
}

func TestNodeRegistryListReturnsSortedNodes(t *testing.T) {
	registry := NewNodeRegistry()

	nodes := []Node{
		testNode("node-03"),
		testNode("node-01"),
		testNode("node-02"),
	}

	for _, node := range nodes {
		if err := registry.Register(node); err != nil {
			t.Fatalf("register %s: %v", node.ID, err)
		}
	}

	got := registry.List()

	if len(got) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(got))
	}

	expected := []string{
		"node-01",
		"node-02",
		"node-03",
	}

	for i, expectedID := range expected {
		if got[i].ID != expectedID {
			t.Fatalf(
				"expected node %d to be %q, got %q",
				i,
				expectedID,
				got[i].ID,
			)
		}
	}
}

func TestNodeRegistryUpdateHeartbeat(t *testing.T) {
	registry := NewNodeRegistry()
	node := testNode("node-01")

	if err := registry.Register(node); err != nil {
		t.Fatalf("register node: %v", err)
	}

	heartbeatAt := node.LastHeartbeat.Add(30 * time.Second)

	if err := registry.UpdateHeartbeat(node.ID, heartbeatAt); err != nil {
		t.Fatalf("update heartbeat: %v", err)
	}

	got, err := registry.Get(node.ID)
	if err != nil {
		t.Fatalf("get node: %v", err)
	}

	if !got.LastHeartbeat.Equal(heartbeatAt) {
		t.Fatalf(
			"expected heartbeat %v, got %v",
			heartbeatAt,
			got.LastHeartbeat,
		)
	}
}

func TestNodeRegistryUpdateHeartbeatUnknownNode(t *testing.T) {
	registry := NewNodeRegistry()

	err := registry.UpdateHeartbeat("missing-node", time.Now())

	if err != ErrNodeNotFound {
		t.Fatalf("expected ErrNodeNotFound, got %v", err)
	}
}

func TestNodeRegistryConcurrentAccess(t *testing.T) {
	registry := NewNodeRegistry()

	const nodeCount = 100

	var wg sync.WaitGroup

	registerErrors := make(chan error, nodeCount)

	for i := 0; i < nodeCount; i++ {
		wg.Add(1)

		go func(index int) {
			defer wg.Done()

			node := testNode(fmt.Sprintf("node-%03d", index))

			if err := registry.Register(node); err != nil {
				registerErrors <- fmt.Errorf(
					"register %s: %w",
					node.ID,
					err,
				)
			}
		}(i)
	}

	wg.Wait()
	close(registerErrors)

	for err := range registerErrors {
		t.Error(err)
	}

	if got := len(registry.List()); got != nodeCount {
		t.Fatalf("expected %d nodes, got %d", nodeCount, got)
	}

	operationErrors := make(chan error, nodeCount*2)

	for i := 0; i < nodeCount; i++ {
		nodeID := fmt.Sprintf("node-%03d", i)

		wg.Add(2)

		go func(id string) {
			defer wg.Done()

			if _, err := registry.Get(id); err != nil {
				operationErrors <- fmt.Errorf(
					"get %s: %w",
					id,
					err,
				)
			}
		}(nodeID)

		go func(id string) {
			defer wg.Done()

			if err := registry.UpdateHeartbeat(id, time.Now()); err != nil {
				operationErrors <- fmt.Errorf(
					"update heartbeat %s: %w",
					id,
					err,
				)
			}
		}(nodeID)
	}

	wg.Wait()
	close(operationErrors)

	for err := range operationErrors {
		t.Error(err)
	}
}
