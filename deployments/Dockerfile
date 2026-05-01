FROM golang:1.24-alpine AS builder

WORKDIR /app

# Izinkan Go mendownload toolchain yang lebih baru jika modul membutuhkannya
ENV GOTOOLCHAIN=auto

# Copy dependency files dulu
COPY go.mod go.sum ./

# Copy semua source code
COPY . .

# Sinkronisasi dan download dependency
RUN go mod tidy


# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app

# Final image — minimal alpine
FROM alpine:3.19

# Install ca-certificates untuk HTTPS call (Google OAuth, Fonnte)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
