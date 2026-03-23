FROM golang:1.23.6 AS builder
WORKDIR /app
COPY . .
RUN go build -o server .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/server /app/server
EXPOSE 8080
CMD ["/app/server"]
