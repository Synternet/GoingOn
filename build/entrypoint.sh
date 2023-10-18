#!/bin/sh

CMD="./goingon"

if [ ! -z "$NATS_URLS" ]; then
  CMD="$CMD --nats-urls $NATS_URLS"
fi

if [ ! -z "$NATS_SUB_NKEY" ]; then
  CMD="$CMD --nats-sub-nkey $NATS_SUB_NKEY"
fi

if [ ! -z "$NATS_RECONNECT_WAIT" ]; then
  CMD="$CMD --nats-reconnect-wait $NATS_RECONNECT_WAIT"
fi

if [ ! -z "$NATS_MAX_RECONNECT" ]; then
  CMD="$CMD --nats-max-reconnect $NATS_MAX_RECONNECT"
fi

if [ ! -z "$NATS_EVENT_LOG_STREAM_SUBJECT" ]; then
  CMD="$CMD --nats-event-log-stream-subject $NATS_EVENT_LOG_STREAM_SUBJECT"
fi

if [ ! -z "$POOL_CONTRACT_ADDRESS" ]; then
  CMD="$CMD --pool-contract-address $POOL_CONTRACT_ADDRESS"
fi

if [ ! -z "$POOL_TOKEN0_TICKER" ]; then
  CMD="$CMD --pool-token0-ticker $POOL_TOKEN0_TICKER"
fi

if [ ! -z "$POOL_TOKEN0_DECIMALS" ]; then
  CMD="$CMD --pool-token0-decimals $POOL_TOKEN0_DECIMALS"
fi

if [ ! -z "$POOL_TOKEN1_TICKER" ]; then
  CMD="$CMD --pool-token1-ticker $POOL_TOKEN1_TICKER"
fi

if [ ! -z "$POOL_TOKEN1_DECIMALS" ]; then
  CMD="$CMD --pool-token1-decimals $POOL_TOKEN1_DECIMALS"
fi

if [ ! -z "$TELEGRAM_CHAT_ID" ]; then
  CMD="$CMD --telegram-chat-id $TELEGRAM_CHAT_ID"
fi

if [ ! -z "$TELEGRAM_BOT_TOKEN" ]; then
  CMD="$CMD --telegram-bot-token $TELEGRAM_BOT_TOKEN"
fi

exec $CMD
