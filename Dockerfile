FROM golang:1.18-alpine
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

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

RUN go generate

# Build the Go app
RUN go build -o /ferrum

COPY --from=ghcr.io/ufoscout/docker-compose-wait:latest /wait /wait

COPY testData ./testData

CMD ["/bin/bash", "-c", "/wait && python ./testData/redis/insert_test_data.py && /ferrum --config ./config_docker_w_redis.json"]
