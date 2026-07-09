FROM golang:1.25 AS builder
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/tcp-server ./services/tcp-server

FROM alpine:latest AS runtime
WORKDIR /app
COPY --from=builder /out/tcp-server ./tcp-server
COPY services/tcp-server/etc ./etc
EXPOSE 7611
ENTRYPOINT ["./tcp-server", "-f", "etc/tcp-server.docker.yaml"]










FROM golang:1.25 AS dlv
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/go-delve/delve/cmd/dlv@latest

FROM builder AS debug-builder
RUN CGO_ENABLED=0 go build -gcflags="all=-N -l" -o /out/tcp-server ./services/tcp-server



FROM golang:1.25 AS debug
WORKDIR /app
COPY --from=debug-builder /out/tcp-server ./tcp-server
COPY --from=dlv /go/bin/dlv /usr/local/bin/dlv
COPY services/tcp-server/etc ./etc
EXPOSE 7611 2348
ENTRYPOINT ["dlv", "--listen=:2348", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "--continue", "/app/tcp-server", "--", "-f", "/app/etc/tcp-server.docker.yaml"]
