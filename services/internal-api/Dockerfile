FROM golang:1.25 AS builder
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/internal-api ./services/internal-api

FROM alpine:latest AS runtime
WORKDIR /app
COPY --from=builder /out/internal-api ./internal-api
COPY services/internal-api/etc ./etc
EXPOSE 8089
ENTRYPOINT ["./internal-api", "-f", "etc/internal-api.docker.yaml"]





FROM golang:1.25 AS dlv
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/go-delve/delve/cmd/dlv@latest

FROM builder AS debug-builder
RUN CGO_ENABLED=0 go build -gcflags="all=-N -l" -o /out/internal-api ./services/internal-api

FROM golang:1.25 AS debug
WORKDIR /app
COPY --from=debug-builder /out/internal-api ./internal-api
COPY --from=dlv /go/bin/dlv /usr/local/bin/dlv
COPY services/internal-api/etc ./etc
EXPOSE 8089 2347
ENTRYPOINT ["dlv", "--listen=:2347", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "--continue", "/app/internal-api", "--", "-f", "/app/etc/internal-api.docker.yaml"]
