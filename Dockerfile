FROM golang:1.18-alpine
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base && apk add --no-cach openssl

RUN mkdir /app
WORKDIR /app

COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

RUN go generate

# Build the Go app
RUN go build -o /ferrum

# Run the executable
CMD [ "/ferrum", "--config", "./config_docker_w_redis.json" ]