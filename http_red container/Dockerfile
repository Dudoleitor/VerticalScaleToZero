FROM golang:1.21.5 AS build-stage
WORKDIR /app
COPY *.go .
COPY go.mod .
COPY go.sum .
RUN CGO_ENABLED=0 GOOS=linux go build -o main

FROM alpine
COPY --from=build-stage /app/main /entrypoint
CMD ["/entrypoint"]
