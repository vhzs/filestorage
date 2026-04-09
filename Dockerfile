FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o filestorage ./cmd/filestorage

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/filestorage .
RUN mkdir -p /app/uploads
EXPOSE 8080
CMD ["./filestorage"]
