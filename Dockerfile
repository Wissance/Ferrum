FROM golang:1.24-alpine3.20
VOLUME /app_data
VOLUME /nginx_cfg

RUN sed -i 's/https/http/' /etc/apk/repositories
RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add --no-cache build-base && apk add --no-cache openssl

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN apk add --update --no-cache python3 && ln -sf python3 /usr/bin/python && apk add py3-pip
RUN apk add py3-setuptools && apk add py3-redis
# TODO(UMV): I need to create dhparam directory in VOLUME, there are no other way or i have not found it yet
# COPY "LICENSE" /nginx_cfg/dhparam/
RUN mkdir -p /nginx_cfg/dhparam && mkdir -p /nginx_cfg/certs && mkdir -p /nginx_cfg/conf.d

#RUN go mod tidy && go generate
# Download all the dependencies
#RUN go get -d -v ./...
#RUN go install -v ./...
RUN go generate

# Build the Go apps
RUN go build -o ferrum
RUN go build -o ferrum-admin ./api/admin/cli

# TODO(SIA) Vulnerability
COPY --from=ghcr.io/ufoscout/docker-compose-wait:latest /wait /wait
COPY tools ./tools

# TODO(UMV): 1. Build config on a Fly (to use props from Env variables)
# TODO(UMV): 2. If we have users, realms and clients do not attempt to insert them

CMD ["/bin/bash", "-c", "/app/tools/docker_app_runner.sh"]
