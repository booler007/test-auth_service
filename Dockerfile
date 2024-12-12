FROM golang:alpine AS builder
LABEL authors="ilya, sewerefd@gmail.com"

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLE=0 COOS=linux GOARCH=amd64 go build -o ./bin/main ./cmd/api

FROM scratch
COPY --from=builder /app/bin/main /app/bin/main

ENTRYPOINT ["/app/bin/main"]
