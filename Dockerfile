FROM golang:1.18-alpine
#VOLUME /app_data

RUN sed -i 's/https/http/' /etc/apk/repositories
RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add --no-cache build-base && apk add --no-cache openssl

RUN apk add --update --no-cache python3 && ln -sf python3 /usr/bin/python
RUN python3 -m ensurepip
RUN pip3 install --no-cache --upgrade pip setuptools

RUN pip install redis

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

RUN go generate

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Build the Go apps
RUN go build -o /ferrum

RUN go build -o ferrum-admin ./api/admin/cli

COPY --from=ghcr.io/ufoscout/docker-compose-wait:latest /wait /wait

COPY testData ./testData

# RUN cp config_docker_w_redis.json /app_data/config_docker_w_redis.json

CMD ["/bin/bash", "-c", "/wait && ./create_wissance_demo_users_docker.sh && /ferrum --config /app/config_docker_w_redis.json"]
