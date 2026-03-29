FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/bcc-skill-gap cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl tzdata && update-ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/bcc-skill-gap .
COPY .env .env

ARG APP_PORT=3007
EXPOSE ${APP_PORT}

ENTRYPOINT [ "./bcc-skill-gap" ]