# Test Containers

Fast-starting database containers optimized for testing.

## Overview

This repository contains Docker containers for databases that are optimized for fast startup in test environments. The key optimization is pre-initializing expensive setup operations during the Docker image build phase, so containers start almost instantly at runtime.

## Available Containers

### PostgreSQL

A PostgreSQL container with pre-initialized data directory.

**Optimization:** The data cluster is created during the Docker build, so at runtime PostgreSQL just starts the server without any initialization.

**Startup time:** ~1 second (vs 3-5 seconds for standard container)

**Connection:** `postgresql://test:test@localhost:5432/test`

See [postgresql/README.md](postgresql/README.md) for details.

### TigerBeetle

A TigerBeetle container with pre-formatted data file.

**Optimization:** The data file (~1GB) is formatted during a one-time setup step, then baked into the image. At runtime, TigerBeetle just starts the server.

**Startup time:** ~5 seconds (avoids the expensive format step at runtime)

**Connection:** Port `3000`, Cluster `0`

See [tigerbeetle/README.md](tigerbeetle/README.md) for details.

## The Pattern

Both containers follow the same optimization pattern:

1. **Two-phase initialization:**
   - **Phase 1 (Build time):** Expensive setup is done once during image build
   - **Phase 2 (Runtime):** Only the server starts, no initialization overhead

2. **Benefits:**
   - Faster test suite execution
   - Consistent initial state for each test
   - No volume mounts needed for basic tests
   - Isolated, ephemeral data perfect for testing
