FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY init.sql /migrations/
CMD ["echo", "Migrations complete"]
