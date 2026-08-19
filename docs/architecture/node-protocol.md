# Clustryn Node Protocol v0.1

## Purpose

This document defines the initial contract between a Clustryn Node Agent and the Control Plane.

A node is a machine or simulated machine capable of joining a Clustryn cluster and reporting its identity and available resources.

---

## Node Registration

### Endpoint

`POST /nodes/register`

### Content Type

`application/json`

### Request

```json
{
  "node_id": "node-01",
  "hostname": "worker-01",
  "cpu_cores": 8,
  "memory_mb": 16384
}