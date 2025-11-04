FROM golang:1.22-alpine as builder

WORKDIR /app/

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app
# ---------------------------------------------------
FROM debian:bullseye-slim as deploy

RUN apt-get update

COPY --from=builder /app/app .

CMD ["./app"]

FROM golang:1.22 as dev
WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
CMD ["air"]