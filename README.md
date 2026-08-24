# MiniTrue

MiniTrue is an early Go time-series ingestion prototype. This repository represents the work completed through 11 September 2026:

```text
synthetic sensor -> HTTP publisher -> ingestion router -> Murmur3 hash ring
```

The publisher generates sensor datapoints or reads them from a serial device. The router accepts `POST /route`, builds a `device_id:metric_name` routing key, and selects a responsible node using a virtual-node consistent hash ring.

## Run

Start the router:

```sh
go run ./cmd/minitrue-router --port 7070 --nodes polaris,sirius,vega
```

Start the simulator in another terminal:

```sh
go run ./cmd/publisher --sim=true --ingest=http://localhost:7070/route
```

The ring implementation and edge cases are covered by `internal/cluster/consistent_hash_test.go`.

This cutoff repository intentionally does not include storage durability, WAL, compression, gossip, query processing, frontend, or deployment services introduced later in the project.
