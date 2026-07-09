FROM golang:1.25 AS builder
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/transform-rpc ./services/rpc-transform






# FROM gcr.io/distroless/static-debian12:nonroot AS runtime
FROM alpine:latest AS runtime
WORKDIR /app
COPY --from=builder /out/transform-rpc ./transform-rpc
COPY services/rpc-transform/etc ./etc
EXPOSE 9001
ENTRYPOINT ["./transform-rpc", "-f", "etc/transform.docker.yaml"]
















FROM golang:1.25 AS dlv
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/go-delve/delve/cmd/dlv@latest

FROM builder AS debug-builder
RUN CGO_ENABLED=0 go build -gcflags="all=-N -l" -o /out/transform-rpc ./services/rpc-transform

FROM golang:1.25 AS debug
WORKDIR /app
COPY --from=debug-builder /out/transform-rpc ./transform-rpc
COPY --from=dlv /go/bin/dlv /usr/local/bin/dlv
COPY services/rpc-transform/etc ./etc
EXPOSE 9001 2345
ENTRYPOINT ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "--continue", "/app/transform-rpc", "--", "-f", "/app/etc/transform.docker.yaml"]
