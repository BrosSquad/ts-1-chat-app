FROM golang:1.16 as builder

COPY . /app

ENV CGO_ENABLED=0

RUN go build -o bin/server main.go


FROM alpine:3.19

COPY --from=builder /app/bin/server /app/server

RUN mkdir -p /var/chat-app/db && mkdir -p /var/logs/chat-app/

ARG PORT=3000

ENV SERVER_ADDR="0.0.0.0:3000"
ENV CONSOLE_LOGGER_LEVEL=trace

EXPOSE ${PORT}

CMD ["/app/server", "-db", "/var/chat-app/db", "-logs", "/var/logs/chat-app", "-level", "${CONSOLE_LOGGER_LEVEL}"]
