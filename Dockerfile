FROM golang:1.16-alpine as builder

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags '-d -w -s' -o main

FROM alpine:3

WORKDIR /app

RUN chown -R 65534:65534 /app

COPY --from=builder --chown=65534:65534 /app/main /app

EXPOSE 8080

# nobody
USER 65534

CMD ["./main"]
