FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY internal/common internal/common
COPY internal/user internal/user
COPY cmd/user cmd/user

WORKDIR /app/internal/user
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service ../../cmd/user/main.go

FROM alpine
WORKDIR /root/
COPY --from=builder /user-service .

CMD ["./user-service"]