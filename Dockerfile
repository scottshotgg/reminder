FROM golang:1.14-alpine3.12 AS builder

ARG TOKEN
ARG FROM
ARG SID
ENV TOKEN ${TOKEN}
ENV FROM ${FROM}
ENV SID ${SID}

WORKDIR /app

COPY . .

RUN go build -o reminder

FROM alpine:3.12

COPY --from=builder /app/reminder /bin/app

CMD ["sh", "-c", "app --sms-provider twilio --from ${FROM} --sid=${SID} --token ${TOKEN} --redis-url redis:6379"]
# CMD ["/app/reminder", "--redis-url", "redis:6379"]