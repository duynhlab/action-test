# Migration image for the docker-db-init job in the build template.
# Build context is db/migrations, so only this directory is visible.
FROM alpine:3.17.0

RUN apk add --no-cache postgresql15-client

WORKDIR /migrations
COPY *.sql ./

ENTRYPOINT ["sh", "-c", "echo 'migrations bundled:' && ls -1 /migrations"]
