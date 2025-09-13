FROM golang:1.24 AS builder

WORKDIR /app

ARG DEBUG=false
ARG DEBUG_PORT=2345

RUN if [ $DEBUG ]; then \
        go install github.com/go-delve/delve/cmd/dlv@latest; \
    fi

COPY go.mod go.sum ./
RUN go mod download

COPY internal/common internal/common
COPY internal/auth internal/auth
COPY cmd/auth cmd/auth

WORKDIR /app/internal/auth

RUN if [ $DEBUG ]; then \
        CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o /auth-service ../../cmd/auth/main.go; \
    else \
        CGO_ENABLED=0 GOOS=linux go build -o /auth-service ../../cmd/auth/main.go; \
    fi

FROM golang AS debug
WORKDIR /root/
COPY --from=builder /auth-service .
COPY --from=builder /go/bin/dlv .
EXPOSE ${DEBUG_PORT}
CMD ["./dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./auth-service"]

FROM gcr.io/distroless/static AS release
WORKDIR /root/
COPY --from=builder /auth-service .
CMD ["./auth-service"]