# Controlled base image for E2.
#
# A mutable tag pointing at this image is the precondition for the "stale base
# layer" hypothesis: the app Dockerfile does FROM <base>:latest, and we rewrite
# that tag underneath it between phases.
#
# UPGRADE=false -> keep the old, vulnerable package set
# UPGRADE=true  -> apk upgrade, which should reduce the CVE count
FROM alpine:3.17.0

ARG UPGRADE=false
RUN if [ "$UPGRADE" = "true" ]; then \
      apk upgrade --no-cache; \
    else \
      echo "keeping the original 3.17.0 package set"; \
    fi

RUN adduser -D -u 10001 app
