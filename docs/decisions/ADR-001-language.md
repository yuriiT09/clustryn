# ADR-001: Go as the Initial Backend Language

## Status

Accepted

## Decision

The initial Control Plane and Node Agent will be written in Go.

## Reason

Clustryn requires networking, concurrency, long-running services, and efficient communication between multiple nodes. Go provides a simple toolchain and strong support for building these types of systems.

## Alternatives Considered

- Rust
- Java
- Python
- C#

## Notes

This decision may be revisited later if another language becomes more suitable for a specific component.