FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .
FROM alpine:latest
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    libstdc++ \
    dumb-init

WORKDIR /app
COPY --from=builder /app/main .
COPY targets.yaml . 
RUN mkdir scans
ENV CHROME_BIN=/usr/bin/chromium-browser
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["./main"]