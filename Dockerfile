FROM golang:1.14-alpine3.12 AS builder

ARG TOKEN
ARG FROM
ARG SID
ARG SMS_PROVIDER
ARG REDIS_URL

ENV TOKEN ${TOKEN}
ENV FROM ${FROM}
ENV SID ${SID}
ENV SMS_PROVIDER ${SMS_PROVIDER}
ENV REDIS_URL ${REDIS_URL}

WORKDIR /app

COPY . .

RUN go build -o reminder

FROM alpine:3.12

COPY --from=builder /app/reminder /bin/app

CMD ["sh", "-c", "app --sms-provider ${SMS_PROVIDER} --from ${FROM} --sid=${SID} --token ${TOKEN} --redis-url ${REDIS_URL}"]
# CMD ["/app/reminder", "--redis-url", "redis:6379"]