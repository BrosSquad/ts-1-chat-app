FROM golang:1.16 as builder

COPY ./backend /app

WORKDIR /app

RUN go build -o ../bin/server ./cmd/server/main.go

FROM golang:1.16

COPY --from=builder /bin/server /bin/server
COPY ./backend/config.example.yml /etc/chat-app/config.yml

RUN mkdir -p /var/chat-app/db && \
    mkdir -p /var/logs/chat-app/ && touch /var/chat-app/db/database.sqlite

ARG PORT=3000

ENV SERVER_ADDR="0.0.0.0:3000"
ENV CONSOLE_LOGGER_LEVEL=trace

EXPOSE ${PORT}

ENTRYPOINT ["server", "-db", "/var/chat-app/db/database.sqlite", "-logs", "/var/logs/chat-app", "-level", "trace", "-config", "/etc/chat-app"]
