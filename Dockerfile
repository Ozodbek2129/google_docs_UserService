FROM golang:1.22.4 AS builder

WORKDIR /app

COPY . .

RUN go mod download

COPY .env .

RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd -a -installsuffix cgo -o ./../myapp .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/myapp .
COPY --from=builder /app/.env .
COPY --from=builder /app/api/email/template.html ./api/email/

EXPOSE 1234
EXPOSE 2345

CMD ["./myapp"]