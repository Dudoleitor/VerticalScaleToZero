FROM golang:1.21.5 AS build-stage
WORKDIR /app
COPY *.go .
COPY go.mod .
RUN CGO_ENABLED=0 GOOS=linux go build -o workload

FROM alpine
COPY --from=build-stage /app/workload /entrypoint
CMD ["/entrypoint"]
