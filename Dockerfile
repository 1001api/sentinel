FROM golang:1.23.2-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/main .
COPY --from=build /app/migrations ./migrations

CMD ["./main"]
