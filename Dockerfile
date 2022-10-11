# syntax=docker/dockerfile:1

## Build
FROM golang:1.18-alpine AS build

WORKDIR /app

# The order is important here, as the go.mod and go.sum files are not changed frequently
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN GOOS=linux GOARCH=amd64 go build -o /app/elevate-sec main.go 


## Deploy
# https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/base-debian10

LABEL   authors="maksym onyshchenko" \
    author-email="placeholder@maxim.run"

WORKDIR /app
COPY --from=build /app/elevate-sec /app/elevate-sec

EXPOSE 9000
USER nonroot:nonroot

ENTRYPOINT ["/app/elevate-sec"]

