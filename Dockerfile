FROM golang:1.18-alpine
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

RUN mkdir /app
WORKDIR /app

COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Build the Go app
RUN go build -o /ferrum

COPY config_docker.json .

ARG SCHEMA
RUN sed -i s/{SCHEMA}/${SCHEMA}/g config_docker.json
ARG HOST
RUN sed -i s/{HOST}/${HOST}/g config_docker.json
ARG PORT
RUN sed -i s/{PORT}/${PORT}/g config_docker.json

ARG REDIS_HOST
RUN sed -i s/{REDIS_HOST}/${REDIS_HOST}/g config_docker.json
ARG REDIS_PORT
RUN sed -i s/{REDIS_PORT}/${REDIS_PORT}/g config_docker.json
ARG REDIS_USER
RUN sed -i s/{REDIS_USER}/${REDIS_USER}/g config_docker.json
ARG REDIS_PASSWORD
RUN sed -i s/{REDIS_PASSWORD}/${REDIS_PASSWORD}/g config_docker.json

ARG NAMESPACE
RUN sed -i s/{NAMESPACE}/${NAMESPACE}/g config_docker.json
ARG DB_NUMBER
RUN sed -i s/{DB_NUMBER}/${DB_NUMBER}/g config_docker.json

# Run the executable
CMD [ "/ferrum", "--config", "./config_docker.json" ]