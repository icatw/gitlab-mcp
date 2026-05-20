FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/mcp-gitlab

FROM alpine:3.22

RUN addgroup -S gitlab-mcp && adduser -S gitlab-mcp -G gitlab-mcp

COPY --from=builder /out/server /usr/local/bin/gitlab-mcp

USER gitlab-mcp
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/gitlab-mcp"]
