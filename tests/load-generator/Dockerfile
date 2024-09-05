FROM golang:1.22.6 AS build-stage
WORKDIR /app
COPY *.go .
COPY go.mod .
COPY go.sum .
RUN CGO_ENABLED=0 GOOS=linux go build -o loadgen

FROM alpine
COPY --from=build-stage /app/loadgen /entrypoint
CMD ["/entrypoint"]
