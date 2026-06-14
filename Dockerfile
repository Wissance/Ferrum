FROM golang:1.24-alpine3.20 AS builder

RUN sed -i 's/https/http/' /etc/apk/repositories
RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add --no-cache build-base && apk add --no-cache openssl
RUN apk add --no-cache ca-certificates tzdata

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go generate

# Build the Go apps
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o ferrum
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o ferrum-admin ./api/admin/cli

FROM alpine:3.20
VOLUME /app_data
VOLUME /nginx_cfg

RUN mkdir -p /nginx_cfg/dhparam && mkdir -p /nginx_cfg/certs && mkdir -p /nginx_cfg/conf.d
# TODO(UMV): I need to create dhparam directory in VOLUME, there are no other way or i have not found it yet
# COPY "LICENSE" /nginx_cfg/dhparam/

RUN apk add --update --no-cache python3 && ln -sf python3 /usr/bin/python && apk add py3-pip
RUN apk add py3-setuptools && apk add py3-redis

RUN addgroup -g 1000 -S wissance && adduser -u 1000 -S ferrum -G wissance
RUN mkdir /app
RUN chown ferrum:wissance /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder --chown=ferrum:wissance --chmod=755 /app/ferrum /app/ferrum
COPY --from=builder --chown=appuser:appgroup --chmod=755 /app/ferrum-admin /app/ferrum-admin
COPY --from=builder --chown=appuser:appgroup --chmod=755 /app/config_docker_w_redis.json /app/config_docker_w_redis.json
COPY --from=builder --chown=ferrum:wissance --chmod=755 /app/tools /app/tools
#TODO(UMV): add keyfile re-generation on every run
COPY --from=builder --chown=ferrum:wissance --chmod=755 /app/keyfile /app/keyfile

WORKDIR /app
USER wissance

CMD ["/bin/sh", "-c", "/app/tools/docker_app_runner.sh"]
