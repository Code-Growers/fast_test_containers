ARG TB_VERSION=0.16.78

FROM ghcr.io/tigerbeetle/tigerbeetle:${TB_VERSION}

# Copy pre-formatted data file
# This file was created by running:
#   docker run --rm --privileged --entrypoint=/bin/sh ghcr.io/tigerbeetle/tigerbeetle:${TB_VERSION} -c \
#     "mkdir -p /data && /tigerbeetle format --cluster=0 --replica=0 --replica-count=1 /data/0_0.tigerbeetle && cat /data/0_0.tigerbeetle" \
#     > 0_0.tigerbeetle
COPY 0_0.tigerbeetle /data/0_0.tigerbeetle

EXPOSE 3000
# Environment variables for configuration
ENV TB_CLUSTER=1
ENV TB_REPLICA=0
ENV TB_ADDRESSES=0.0.0.0:3000
ENV TB_DATA_FILE=/data/0_0.tigerbeetle
# Use smaller cache for tests to reduce memory footprint
ENV TB_CACHE_GRID=256MiB

# Start TigerBeetle server directly with pre-formatted data
# Note: Using --cache-grid for lower memory usage in test environments
ENTRYPOINT ["/tigerbeetle"]

STOPSIGNAL SIGINT

CMD ["start", "--addresses=0.0.0.0:3000", "--cache-grid=256MiB", "/data/0_0.tigerbeetle"]
