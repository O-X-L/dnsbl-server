FROM golang:1.24.5-alpine AS builder

# BUILD: docker build --no-cache -t dnsbl-server -f Dockerfile .
# TEST: docker run -it --rm --name dnsbl-server -p 5353:5353 dnsbl-server:latest /usr/local/bin/dnsbl-server -help
# RUN: docker run -d --name dnsbl-server --restart always -p 53:5353/udp -v $(pwd)/config.yml:/app/config.yml dnsbl-server:latest /usr/local/bin/dnsbl-server -config /app/config.yml -domain test.at

ENV REPO_URL=https://github.com/O-X-L/dnsbl-server/archive/refs/heads/latest.tar.gz

RUN set -eux; \
        wget -O /tmp/dnsbl-server.tar.gz "$REPO_URL"; \
        tar -xzf /tmp/dnsbl-server.tar.gz -C /tmp/; \
        rm /tmp/dnsbl-server.tar.gz; \
        cd /tmp/dnsbl-server-latest; \
        go mod tidy; \
        go build -o /tmp/dnsbl-server /tmp/dnsbl-server-latest/src/cmd/

FROM alpine:latest

RUN set -eux; \
        addgroup --gid 7005 --system dnsbl; \
        adduser \
                --disabled-password \
                --home /app \
                --ingroup dnsbl \
                --no-create-home \
                --system \
                --uid 7005 \
                dnsbl \
        ; \
        mkdir /app; \
        chown dnsbl:dnsbl /app

COPY --from=builder /tmp/dnsbl-server /usr/local/bin

USER dnsbl
WORKDIR /app
