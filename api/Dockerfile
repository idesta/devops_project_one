# api/Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY main.go ./
RUN apk add --no-cache git build-base
RUN go build -o /app .

FROM alpine:3.18
COPY --from=builder /app /app
EXPOSE 8080
ENV PORT=8080
CMD ["/app"]