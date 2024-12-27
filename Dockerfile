FROM golang:1.23.2-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags '-w -s' -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/main .
COPY --from=build /app/migrations ./internal/migrations
COPY --from=build /app/views ./views
COPY --from=build /app/ipdb ./internal/ipdb

CMD ["./main"]
