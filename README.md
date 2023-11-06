# GoingOn

GoingOn lets you filter and decode Uniswap pool smart contract swaps from Ethereum transaction receipt event log.

## Live Demo

[Uniswap NOIA-ETH pool](https://v2.info.uniswap.org/pair/0xb8a1a865e4405281311c5bc0f90c240498472d3e) demo available at https://t.me/+KUUWjzsV6v04ZDU0.

## Usage

1. Compile code.
```
make build
```

2. Run executable.
```
./goingon [flags]
```

## Flags

| Flag                          | Description                                                       |
| ----------------------------- | ----------------------------------------------------------------- |
| nats-urls                     | NATS servers URLs (comma separated)                               |
| nats-sub-nkey                 | NATS user credentials NKey string                                 |
| nats-reconnect-wait           | NATS reconnect wait duration                                      |
| nats-max-reconnect            | NATS max reconnect attempts count                                 |
| nats-event-log-stream-subject | NATS event log stream subject                                     |
| pool-contract-address         | Pool contract address                                             |
| pool-token0-ticker            | Ticker of first pool token                                        |
| pool-token0-decimals          | Decimals of first pool token                                      |
| pool-token1-ticker            | Ticker of second pool token                                       |
| pool-token1-decimals          | Decimals of second pool token                                     |
| telegram-chat-id              | Telegram chat ID                                                  |
| telegram-bot-token            | Telegram bot token                                                |

- `nats-*`. NATS. `nats-sub-nkey` must be provided. Uses Syntropy Data Layer to get Ethereum transactions event log. See [Data Layer Quick Start](https://docs.syntropynet.com/build/) to learn more.
- `pool-*`. Uniswap Pool. Default flags is set to [Uniswap ETH-USDT Pool](https://v2.info.uniswap.org/pair/0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852).
- `telegram-*`. Telegram. Flags must be provided. See [Telegram From BotFather to 'Hello World'](https://core.telegram.org/bots/tutorial) to learn more.

## Docker

1. Build image.
```
docker build -f ./docker/Dockerfile -t goingon .
```

2. Run container with passed environment variables.
```
docker run -it --rm --env-file=.env going-on
```

Note: [Flags](#flags) can be passed as environment variables.
Environment variables are all caps flags separated with underscore. See `./docker/entrypoint.sh`.

## Contributing

We welcome contributions from the community. Whether it's a bug report, a new feature, or a code fix, your input is valued and appreciated.

## Syntropy

If you have any questions, ideas, or simply want to connect with us, we encourage you to reach out through any of the following channels:

- **Discord**: Join our vibrant community on Discord at [https://discord.com/invite/jqZur5S3KZ](https://discord.com/invite/jqZur5S3KZ). Engage in discussions, seek assistance, and collaborate with like-minded individuals.
- **Telegram**: Connect with us on Telegram at [https://t.me/SyntropyNet](https://t.me/SyntropyNet). Stay updated with the latest news, announcements, and interact with our team members and community.
- **Email**: If you prefer email communication, feel free to reach out to us at devrel@syntropynet.com. We're here to address your inquiries, provide support, and explore collaboration opportunities.
