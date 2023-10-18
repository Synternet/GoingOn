# Build Container
FROM golang:1.18 AS build

ARG VERSION
ARG NATS_URLS NATS_SUB_CREDS NATS_SUB_JWT NATS_SUB_NKEY NATS_RECONNECT_WAIT NATS_MAX_RECONNECT NATS_EVENT_LOG_STREAM_SUBJECT
ARG POOL_CONTRACT_ADDRESS POOL_TOKEN0_TICKER POOL_TOKEN0_DECIMALS POOL_TOKEN1_TICKER POOL_TOKEN1_DECIMALS
ARG TELEGRAM_CHAT_ID TELEGRAM_BOT_TOKEN

COPY . /home/src
WORKDIR /home/src

RUN CGO_ENABLED=0 go build -v -o . -installsuffix cgo -ldflags="-w -s" ./...

# Runtime Container
FROM alpine

ENV PATH="/app:${PATH}"

RUN apk add --update --no-cache \
    ca-certificates \
    bash

WORKDIR /home/app
COPY --from=build /home/src/goingon .
COPY --from=build /home/src/build/entrypoint.sh /home/app
RUN chmod +x /home/app/entrypoint.sh

ENTRYPOINT ["/home/app/entrypoint.sh"]
