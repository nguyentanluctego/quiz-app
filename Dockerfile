# Build stage
FROM --platform=linux/amd64 golang:1.23.3-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o quiz-app ./cmd/server

# Production stage
FROM --platform=linux/amd64 alpine:latest 

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/quiz-app .
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./quiz-app"] 