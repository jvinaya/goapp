FROM golang:1.17-alpine3.15 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.15
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY SQL/migrations ./db/migrations

EXPOSE 14060
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
