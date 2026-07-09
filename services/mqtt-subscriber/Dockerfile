FROM golang:1.25 AS builder
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/mqtt-subscriber ./services/mqtt-subscriber

FROM alpine:latest AS runtime
WORKDIR /app
COPY --from=builder /out/mqtt-subscriber ./mqtt-subscriber
COPY services/mqtt-subscriber/etc ./etc
ENTRYPOINT ["./mqtt-subscriber", "-f", "etc/mqtt-subscriber.docker.yaml"]






FROM golang:1.25 AS dlv
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/go-delve/delve/cmd/dlv@latest

FROM builder AS debug-builder
RUN CGO_ENABLED=0 go build -gcflags="all=-N -l" -o /out/mqtt-subscriber ./services/mqtt-subscriber

FROM golang:1.25 AS debug
WORKDIR /app
COPY --from=debug-builder /out/mqtt-subscriber ./mqtt-subscriber
COPY --from=dlv /go/bin/dlv /usr/local/bin/dlv
COPY services/mqtt-subscriber/etc ./etc
ENTRYPOINT ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "--continue", "/app/mqtt-subscriber", "--", "-f", "/app/etc/mqtt-subscriber.docker.yaml"]
