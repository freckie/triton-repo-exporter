FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o triton-repo-exporter .

# ---

FROM alpine

WORKDIR /app
COPY --from=builder /app/triton-repo-exporter .

RUN apk --no-cache add ca-certificates tzdata

ENTRYPOINT [ "/app/triton-repo-exporter" ]