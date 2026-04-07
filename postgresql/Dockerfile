# syntax=docker/dockerfile:1

ARG PG_VERSION=17.4

# Stage 1: Initialize PostgreSQL data directory
FROM postgres:${PG_VERSION}-alpine AS initializer

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_DB=postgres
ENV PGDATA=/var/lib/postgresql/data

# Run the postgres entrypoint in initialization-only mode to create the data directory
# This performs all the expensive setup (creating cluster, system tables, etc.)
RUN docker-entrypoint.sh postgres --version \
    && su - postgres -c "pg_ctl initdb -D /var/lib/postgresql/data" \
    && echo "host all all all scram-sha-256" >> /var/lib/postgresql/data/pg_hba.conf \
    && echo "listen_addresses='*'" >> /var/lib/postgresql/data/postgresql.conf \
    && echo "fsync = off" >> /var/lib/postgresql/data/postgresql.conf \
    && echo "synchronous_commit = off" >> /var/lib/postgresql/data/postgresql.conf \
    && echo "full_page_writes = off" >> /var/lib/postgresql/data/postgresql.conf \
    && echo "wal_level = minimal" >> /var/lib/postgresql/data/postgresql.conf \
    && echo "max_wal_senders = 0" >> /var/lib/postgresql/data/postgresql.conf \
    && echo "autovacuum = off" >> /var/lib/postgresql/data/postgresql.conf \
    && echo "log_statement = 'none'" >> /var/lib/postgresql/data/postgresql.conf \
    && echo "log_min_duration_statement = -1" >> /var/lib/postgresql/data/postgresql.conf \
    && mkdir -p /var/run/postgresql && chown postgres:postgres /var/run/postgresql \
    && su - postgres -c "pg_ctl start -D /var/lib/postgresql/data" \
    && sleep 2 \
    && su - postgres -c "psql -c \"ALTER USER postgres WITH PASSWORD 'postgres';\"" \
    && su - postgres -c "pg_ctl stop -D /var/lib/postgresql/data -m fast" \
    && rm -rf /var/lib/postgresql/data/postmaster.pid \
    && rm -rf /var/lib/postgresql/data/log/*

# Stage 2: Final minimal image with pre-initialized data
ARG PG_VERSION
FROM postgres:${PG_VERSION}-alpine

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_DB=postgres
ENV PGDATA=/var/lib/postgresql/data

# Copy the pre-initialized data directory from stage 1
COPY --from=initializer /var/lib/postgresql/data /var/lib/postgresql/data

# Ensure proper ownership
RUN chown -R postgres:postgres /var/lib/postgresql/data \
    && chmod 700 /var/lib/postgresql/data \
    && mkdir -p /var/run/postgresql \
    && chown postgres:postgres /var/run/postgresql \
    && chmod 2777 /var/run/postgresql

# Skip initialization since data is already prepared
ENV POSTGRES_INITDB_ARGS=""
ENV PGUSER=postgres

EXPOSE 5432

# Start postgres directly - no initialization needed!
USER postgres
CMD ["postgres", "-c", "listen_addresses=*"]
