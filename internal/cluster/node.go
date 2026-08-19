package cluster

import "time"

type NodeStatus string

const (
	NodeStatusHealthy NodeStatus = "HEALTHY"
	NodeStatusSuspect NodeStatus = "SUSPECT"
	NodeStatusDead    NodeStatus = "DEAD"
)

// Node represents a machine participating in a Clustryn cluster.
type Node struct {
	ID            string
	Hostname      string
	CPUCores      uint16
	MemoryMB      uint64
	Status        NodeStatus
	RegisteredAt  time.Time
	LastHeartbeat time.Time
}
