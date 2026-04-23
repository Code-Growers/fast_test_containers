ARG TB_VERSION=0.16.78

FROM ghcr.io/tigerbeetle/tigerbeetle:${TB_VERSION}

# Copy pre-formatted data file
# This file was created by running:
#   docker run --rm --privileged --entrypoint=/bin/sh ghcr.io/tigerbeetle/tigerbeetle:${TB_VERSION} -c \
#     "mkdir -p /data && /tigerbeetle format --cluster=1 --replica=0 --replica-count=1 /data/0_0.tigerbeetle"
#   docker cp <container>:/data/0_0.tigerbeetle ./tigerbeetle/
COPY 0_0.tigerbeetle /data/0_0.tigerbeetle

EXPOSE 3000
# Environment variables for configuration
ENV TB_CLUSTER=1
ENV TB_REPLICA=0
ENV TB_ADDRESSES=0.0.0.0:3000
ENV TB_DATA_FILE=/data/0_0.tigerbeetle

# Use tini as init so TigerBeetle does not run as PID 1.
#gWithout tini, signals like SIGTERM are ignored by the kernel for PID 1,
# causing docker stop to wait the full timeout (10s) before SIGKILL.
ENTRYPOINT ["tini", "--", "/tigerbeetle"]

# Use --development for minimal memory footprint and test-optimized settings.
# This reduces per-instance RSS from ~2.3GiB to ~1.4GiB, allowing many
# parallel containers without swapping or OOMs.
CMD ["start", "--addresses=0.0.0.0:3000", "--development", "/data/0_0.tigerbeetle"]
