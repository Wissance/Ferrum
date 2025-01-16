FROM golang:1.22-alpine
VOLUME /app_data
VOLUME /nginx_cfg

RUN sed -i 's/https/http/' /etc/apk/repositories
RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add --no-cache build-base && apk add --no-cache openssl

RUN apk add --update --no-cache python3 && ln -sf python3 /usr/bin/python && apk add py3-pip
RUN apk add py3-setuptools && apk add py3-redis

RUN mkdir /app
WORKDIR /app

COPY api ./api
COPY application ./application
COPY certs ./certs
COPY config ./config
COPY data ./data
COPY dto ./dto
COPY errors ./errors
COPY globals ./globals
COPY logging ./logging
COPY managers ./managers
COPY services ./services
COPY utils ./utils
COPY "go.mod" ./"go.mod"
COPY "go.sum" ./"go.sum"
COPY keyfile ./keyfile
COPY "main.go" ./"main.go"
COPY "config_docker_w_redis.json" ./"config_docker_w_redis.json"
COPY tools/"create_wissance_demo_users_docker.sh" ./"create_wissance_demo_users_docker.sh"
COPY tools/"docker_app_runner.sh" ./"docker_app_runner.sh"
COPY docs/nginx/"nginx_docker.conf" /nginx_cfg/"nginx.conf"
COPY swagger ./swagger
# TODO(UMV): I need to create dhparam directory in VOLUME, there are no other way or i have not found it yet
# COPY "LICENSE" /nginx_cfg/dhparam/
RUN mkdir -p /nginx_cfg/dhparam && mkdir -p /nginx_cfg/certs && mkdir -p /nginx_cfg/conf.d

RUN go mod tidy && go generate
# Download all the dependencies
RUN go get -d -v ./...
RUN go install -v ./...

# Build the Go apps
RUN go build -o ferrum
RUN go build -o ferrum-admin ./api/admin/cli

# TODO(SIA) Vulnerability
COPY --from=ghcr.io/ufoscout/docker-compose-wait:latest /wait /wait

COPY testData ./testData
COPY tools ./tools

# TODO(UMV): 1. Build config on a Fly (to use props from Env variables)

# TODO(UMV): 2. If we have users, realms and clients do not attempt to insert them

CMD ["/bin/bash", "-c", "./docker_app_runner.sh"]
