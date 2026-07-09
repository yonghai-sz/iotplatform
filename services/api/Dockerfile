FROM golang:1.25 AS builder
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/platform-api ./services/api




# FROM gcr.io/distroless/static-debian12:nonroot AS runtime
FROM alpine:latest AS runtime
# USER nobody
WORKDIR /app
COPY --from=builder /out/platform-api ./platform-api
COPY services/api/etc ./etc
EXPOSE 8080
ENTRYPOINT ["./platform-api", "-f", "etc/platform-api.docker.yaml"]









   











FROM golang:1.25 AS dlv
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/go-delve/delve/cmd/dlv@latest

FROM builder AS debug-builder
RUN CGO_ENABLED=0 go build -gcflags="all=-N -l" -o /out/platform-api ./services/api
  
FROM golang:1.25 AS debug
WORKDIR /app
COPY --from=debug-builder /out/platform-api ./platform-api
COPY --from=dlv /go/bin/dlv /usr/local/bin/dlv
COPY services/api/etc ./etc
EXPOSE 8080 2345
ENTRYPOINT ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "--continue", "/app/platform-api", "--", "-f", "/app/etc/platform-api.docker.yaml"]
