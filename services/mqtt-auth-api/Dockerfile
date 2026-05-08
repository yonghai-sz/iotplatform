FROM golang:1.24 AS builder
WORKDIR /src
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/mqtt-auth-api ./services/mqtt-auth-api

FROM alpine:latest AS runtime
WORKDIR /app
COPY --from=builder /out/mqtt-auth-api ./mqtt-auth-api
COPY services/mqtt-auth-api/etc ./etc
EXPOSE 8089
ENTRYPOINT ["./mqtt-auth-api", "-f", "etc/mqtt-auth-api.docker.yaml"]

FROM golang:1.24 AS dlv
WORKDIR /src
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/go-delve/delve/cmd/dlv@latest

FROM builder AS debug-builder
RUN CGO_ENABLED=0 go build -gcflags="all=-N -l" -o /out/mqtt-auth-api ./services/mqtt-auth-api

FROM golang:1.24 AS debug
WORKDIR /app
COPY --from=debug-builder /out/mqtt-auth-api ./mqtt-auth-api
COPY --from=dlv /go/bin/dlv /usr/local/bin/dlv
COPY services/mqtt-auth-api/etc ./etc
EXPOSE 8089 2347
ENTRYPOINT ["dlv", "--listen=:2347", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "--continue", "/app/mqtt-auth-api", "--", "-f", "/app/etc/mqtt-auth-api.docker.yaml"]
