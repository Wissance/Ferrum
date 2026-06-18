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

RUN apk add --update --no-cache python3 && ln -sf python3 /usr/bin/python && apk add py3-pip
RUN python ./tools/lf_fixer.py --i=./tools --o=./tools --sel=*.sh

FROM alpine:3.20
VOLUME /app_data

#RUN apk add --update --no-cache python3 && ln -sf python3 /usr/bin/python && apk add py3-pip
#RUN apk add py3-setuptools && apk add py3-redis
RUN apk add --no-cache redis

RUN addgroup -g 1000 -S wissance && adduser -u 1000 -S ferrum -G wissance
RUN mkdir /app
RUN chown ferrum:wissance /app
WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder --chown=ferrum:wissance --chmod=755 /app/ferrum ./ferrum
COPY --from=builder --chown=appuser:appgroup --chmod=755 /app/ferrum-admin ./ferrum-admin
#TODO(UMV): pass desired config via env, and make it RO ??
COPY --from=builder --chown=appuser:appgroup --chmod=755 /app/config_docker_w_redis.json ./config_docker_w_redis.json
COPY --from=builder --chown=ferrum:wissance --chmod=755 /app/tools/*.sh ./tools/
#TODO(UMV): add keyfile re-generation on every run (via go generate at builder)
COPY --from=builder --chown=ferrum:wissance --chmod=755 /app/keyfile ./keyfile
COPY --from=builder --chown=ferrum:wissance --chmod=755 /app/swagger ./swagger
COPY --from=builder --chown=ferrum:wissance --chmod=755 /app/certs ./certs

USER wissance

CMD ["/bin/sh", "-c", "/app/tools/docker_app_runner.sh"]
