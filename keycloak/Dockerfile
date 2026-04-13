ARG KC_VERSION=26.6.0

# ── Stage 1: Quarkus augmentation (build-time compilation) ───────────────────
# Runs `kc.sh build` which performs Quarkus AOT augmentation: generates
# optimized bytecode, locks in build-time config, and eliminates the
# re-augmentation penalty on every container start.
FROM quay.io/keycloak/keycloak:${KC_VERSION} AS builder

ENV KC_HEALTH_ENABLED=false
ENV KC_METRICS_ENABLED=false
ENV KC_DB=dev-file

WORKDIR /opt/keycloak

RUN /opt/keycloak/bin/kc.sh build

# ── Stage 2: Database pre-initialization ─────────────────────────────────────
# The single largest startup cost on a fresh container is the first-boot
# database initialization: Keycloak must create all H2 schema objects and
# bootstrap the master realm (admin user, default clients, roles, scopes).
# We start the server in the background, read its logs until "started in"
# appears, then send SIGTERM and wait for a clean shutdown — leaving a fully
# initialized H2 data directory baked into the image layer.
FROM quay.io/keycloak/keycloak:${KC_VERSION} AS initializer

COPY --from=builder /opt/keycloak/ /opt/keycloak/

ENV KC_HTTP_ENABLED=true
ENV KC_BOOTSTRAP_ADMIN_USERNAME=admin
ENV KC_BOOTSTRAP_ADMIN_PASSWORD=admin
ENV KC_HOSTNAME_STRICT=false
ENV KC_DB=dev-file

RUN /bin/bash -c ' \
    /opt/keycloak/bin/kc.sh start --optimized > /tmp/kc.log 2>&1 & \
    KC_PID=$!; \
    echo "Keycloak PID: $KC_PID"; \
    timeout=240; elapsed=0; \
    while [ $elapsed -lt $timeout ]; do \
        if grep -q "started in" /tmp/kc.log 2>/dev/null; then \
            echo "Keycloak initialized — shutting down."; \
            kill $KC_PID; \
            wait $KC_PID 2>/dev/null || true; \
            rm -f /tmp/kc.log; \
            echo "Pre-initialization complete."; \
            exit 0; \
        fi; \
        if ! kill -0 $KC_PID 2>/dev/null; then \
            echo "ERROR: Keycloak process died unexpectedly:" >&2; \
            cat /tmp/kc.log >&2; \
            exit 1; \
        fi; \
        sleep 2; \
        elapsed=$((elapsed + 2)); \
    done; \
    echo "ERROR: Keycloak did not start within ${timeout}s. Log:" >&2; \
    cat /tmp/kc.log >&2; \
    kill $KC_PID 2>/dev/null; \
    exit 1 \
'

# ── Stage 3: Final minimal image ─────────────────────────────────────────────
FROM quay.io/keycloak/keycloak:${KC_VERSION}

# Copy augmented server (from kc.sh build)
COPY --from=builder /opt/keycloak/ /opt/keycloak/

# Copy pre-initialized H2 database (master realm already bootstrapped)
COPY --from=initializer /opt/keycloak/data/ /opt/keycloak/data/

ENV KC_HTTP_ENABLED=true
ENV KC_BOOTSTRAP_ADMIN_USERNAME=admin
ENV KC_BOOTSTRAP_ADMIN_PASSWORD=admin
ENV KC_HOSTNAME_STRICT=false
ENV KC_SERVER_ASYNC_BOOTSTRAP=true

ENTRYPOINT ["/opt/keycloak/bin/kc.sh", "start", "--optimized"]
