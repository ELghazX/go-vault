FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o go-vault ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite

WORKDIR /app

COPY --from=builder /app/go-vault .

COPY --from=builder /app/static ./static
COPY --from=builder /app/migrations ./migrations

RUN mkdir -p uploads

RUN addgroup -g 1001 -S appgroup && \
  adduser -u 1001 -S appuser -G appgroup

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

CMD ["./go-vault"]
