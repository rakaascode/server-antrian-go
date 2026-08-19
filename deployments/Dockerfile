FROM golang:1.25-alpine AS builder

WORKDIR /app

# 1. Copy file dependency terlebih dahulu
COPY go.mod go.sum ./

# 2. Download dependency (akan di-cache oleh Docker selama go.mod & go.sum tidak berubah)
RUN go mod download

# 3. Copy seluruh source code
COPY . .

# 4. Build binary utama dan background worker
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/app && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main_worker ./cmd/worker

# Final image — minimal alpine
FROM alpine:3.19

# Install ca-certificates untuk HTTPS call (Google OAuth, Fonnte)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/main_worker .

EXPOSE 8080

CMD ["./main"]
