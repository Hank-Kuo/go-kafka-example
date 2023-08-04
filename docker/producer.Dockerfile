FROM golang:1.19-alpine as builder

WORKDIR /app

COPY ./ /app

RUN go mod download

RUN CGO_ENABLED=0 go build -o bin/ cmd/producer/main.go

EXPOSE 8000

VOLUME [ "./assets" ]
ENTRYPOINT ["./bin/main"]