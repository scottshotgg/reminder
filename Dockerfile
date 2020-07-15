FROM golang:1.14-alpine3.12 AS builder

WORKDIR /app

COPY . .

RUN go build -o reminder

FROM alpine:3.12

COPY --from=builder /app/reminder /bin/reminder

CMD ["reminder"]