FROM golang:1.22.1-alpine3.19 AS builder

COPY . /github.com/gomscourse/chat-server/source/
WORKDIR /github.com/gomscourse/chat-server/source/

RUN go mod download
RUN go build -o ./bin/chat_server cmd/main.go

FROM alpine:latest

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/*

WORKDIR /root/
COPY --from=builder /github.com/gomscourse/chat-server/source/bin/chat_server .
COPY --from=builder /github.com/gomscourse/chat-server/source/entrypoint.sh .
COPY --from=builder /github.com/gomscourse/chat-server/source/migrations ./migrations

ADD https://github.com/pressly/goose/releases/download/v3.14.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose