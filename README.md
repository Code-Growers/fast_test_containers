# Test Containers

Fast-starting database containers optimized for testing.

## Overview

This repository contains Docker containers for databases that are optimized for fast startup in test environments. The key optimization is pre-initializing expensive setup operations during the Docker image build phase, so containers start almost instantly at runtime.

## Available Containers

### PostgreSQL

A PostgreSQL container with pre-initialized data directory.

**Optimization:** The data cluster is created during the Docker build, so at runtime PostgreSQL just starts the server without any initialization.

**Startup time:** ~1 second (vs 3-5 seconds for standard container)

**Connection:** `postgresql://postgres:postgres@localhost:5432/postgres`

### TigerBeetle

A TigerBeetle container with pre-formatted data file.

**Optimizations:**
- The data file (~1GB) is formatted during a one-time setup step, then baked into the image.
- Uses `--development` mode for minimal memory footprint (~1.4GB vs ~2.3GB).
- Uses `tini` as the init process so `docker stop` takes ~0.3s instead of the default 10s timeout (the upstream binary ignores signals when running as PID 1).

**Startup time:** ~1 second

**Shutdown time:** ~0.3 seconds

**Connection:** Port `3000`, Cluster `1`

**Required runtime flags:** `--security-opt seccomp=unconfined` (TigerBeetle uses io_uring, which Docker blocks by default).

### Keycloak

A Keycloak container with pre-initialized H2 database.

**Optimization:** The H2 database is initialized during the Docker build, so at runtime Keycloak just starts the server without any initialization.

**Startup time:** ~2 seconds (vs ~30 seconds for standard container)

**Connection:** Port `8080`, Realm `master`, User `admin`, Password `admin`

## The Pattern

All containers follow the same optimization pattern:

1. **Two-phase initialization:**
   - **Phase 1 (Build time):** Expensive setup is done once during image build
   - **Phase 2 (Runtime):** Only the server starts, no initialization overhead

2. **Benefits:**
   - Faster test suite execution
   - Consistent initial state for each test
   - No volume mounts needed for basic tests
   - Isolated, ephemeral data perfect for testing

## Usage Examples

### Go (testcontainers-go)

See [`examples/go/tigerbeetle_test.go`](examples/go/tigerbeetle_test.go) for a working example that starts multiple TigerBeetle instances in parallel using `testcontainers-go`.

Key points when using this image programmatically:

- Always set `--security-opt seccomp=unconfined` in the host config.
- Wait for the log `"listening on"` rather than just an open port, because TigerBeetle opens the port before it is fully ready.
- Terminate containers normally; because the image uses `tini`, shutdown is instant.
