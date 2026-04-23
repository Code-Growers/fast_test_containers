package tigerbeetle_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TigerBeetleImage is the pre-built fast-starting test container image.
const TigerBeetleImage = "ghcr.io/code-growers/fast_test_containers/tigerbeetle-test:0.16.78"

// startTigerBeetle starts a single TigerBeetle container with the security
// settings required by the upstream binary (io_uring needs seccomp=unconfined).
func startTigerBeetle(ctx context.Context, t *testing.T) (testcontainers.Container, string, error) {
	t.Helper()

	req := testcontainers.ContainerRequest{
		Image:        TigerBeetleImage,
		ExposedPorts: []string{"3000/tcp"},
		// TigerBeetle uses io_uring which is blocked by Docker's default seccomp
		// profile. Without this the container will fail immediately with
		// "error: PermissionDenied".
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.SecurityOpt = []string{"seccomp=unconfined"}
		},
		// Wait until the replica reports it is listening. This is much faster
		// than a port-based wait because the port is open before TB is ready.
		WaitingFor: wait.ForLog("listening on").WithStartupTimeout(10 * time.Second),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("start tigerbeetle: %w", err)
	}

	mappedPort, err := c.MappedPort(ctx, "3000/tcp")
	if err != nil {
		_ = c.Terminate(ctx)
		return nil, "", fmt.Errorf("get mapped port: %w", err)
	}

	return c, mappedPort.Port(), nil
}

func TestTigerBeetle_SingleInstance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	c, port, err := startTigerBeetle(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := c.Terminate(ctx); err != nil {
			t.Logf("terminate container: %v", err)
		}
	}()

	t.Logf("TigerBeetle ready on port %s", port)
}

func TestTigerBeetle_ParallelInstances(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	// Start many containers in parallel. Because the image uses tini as the
	// init process and --development for minimal memory usage, each instance
	// starts in ~1s and stops in ~0.3s instead of hitting the 10s Docker stop
	// timeout.
	for i := range 5 {
		t.Run(fmt.Sprintf("replica-%d", i), func(t *testing.T) {
			t.Parallel()

			c, port, err := startTigerBeetle(ctx, t)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := c.Terminate(ctx); err != nil {
					t.Logf("terminate container: %v", err)
				}
			}()

			t.Logf("TigerBeetle replica %d ready on port %s", i, port)
		})
	}
}
