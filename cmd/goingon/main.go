package main

import (
	"context"
	"flag"
	"log"
	"math/big"
	"time"

	svcnats "github.com/synternet/goingon/pkg/nats"

	"github.com/synternet/goingon/internal/service"

	nats "github.com/nats-io/nats.go"
)

func main() {
	flagNatsUrls := flag.String("nats-urls", "nats://34.107.87.29 ", "NATS server URLs (separated by comma)")
	flagUserCredsSeedSub := flag.String("nats-sub-nkey", "", "NATS User credentials NKey string")
	flagNatsReconnectWait := flag.Duration("nats-reconnect-wait", 10*time.Second, "NATS reconnect wait duration")
	flagNatsMaxReconnects := flag.Int("nats-max-reconnect", 500, "NATS max reconnect attempts count")
	flagNatsTxLogEventsStreamSubject := flag.String("nats-event-log-stream-subject", "synternet.ethereum.log-event", "NATS event log stream subject")
	flagPoolContractAddress := flag.String("pool-contract-address", "0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852", "Pool contract address")
	flagPoolToken0Ticker := flag.String("pool-token0-ticker", "ETH", "Ticker of first pool token")
	flagPoolToken0DecimalsNum := flag.Int64("pool-token0-decimals", 18, "Decimals of first pool token")
	flagPoolToken1Ticker := flag.String("pool-token1-ticker", "USDT", "Ticker of second pool token")
	flagPoolToken1DecimalsNum := flag.Int64("pool-token1-decimals", 6, "Decimals of second pool token")
	flagTelegramChatID := flag.Int64("telegram-chat-id", 0, "Telegram chat ID")
	flagTelegramBotToken := flag.String("telegram-bot-token", "", "Telegram bot token")

	flag.Parse()

	checkRequiredFlags(
		flagPoolContractAddress,
		flagPoolToken0Ticker,
		flagPoolToken0DecimalsNum,
		flagPoolToken1Ticker,
		flagPoolToken1DecimalsNum,
		flagTelegramChatID,
		flagTelegramBotToken,
	)

	token0Decimals := new(big.Int).Exp(big.NewInt(10), big.NewInt(*flagPoolToken0DecimalsNum), nil)
	token1Decimals := new(big.Int).Exp(big.NewInt(10), big.NewInt(*flagPoolToken1DecimalsNum), nil)

	optsSub := []nats.Option{}

	flagUserCredsJWTSub, err := svcnats.CreateAppJwt(*flagUserCredsSeedSub)
	if err != nil {
		log.Fatalf("failed to create JWT: %v", err)
	}
	optsSub = append(optsSub, nats.UserJWTAndSeed(flagUserCredsJWTSub, *flagUserCredsSeedSub))

	optsSub = append(optsSub, nats.MaxReconnects(*flagNatsMaxReconnects))
	optsSub = append(optsSub, nats.ReconnectWait(*flagNatsReconnectWait))

	svcnSub := svcnats.MustConnect(
		svcnats.Config{
			URI:  *flagNatsUrls,
			Opts: optsSub,
		})
	log.Println("Service connected to NATS")

	cfg := service.Config{
		TelegramBotToken: *flagTelegramBotToken,
		TelegramChatID:   *flagTelegramChatID,
		ContractAddress:  *flagPoolContractAddress,
		Token0: service.Token{
			Ticker:   *flagPoolToken0Ticker,
			Decimals: token0Decimals,
		},
		Token1: service.Token{
			Ticker:   *flagPoolToken1Ticker,
			Decimals: token1Decimals,
		},
	}
	sSub := service.NewService(svcnSub, context.Background(), cfg)

	svcnSub.AddHandler(*flagNatsTxLogEventsStreamSubject, sSub.ProcessTxLogEventFromStream)

	sSub.Serve()
}

func checkRequiredFlags(
	poolContractAddress *string,
	poolToken0Ticker *string,
	poolToken0DecimalsNum *int64,
	poolToken1Ticker *string,
	poolToken1DecimalsNum *int64,
	telegramChatID *int64,
	telegramBotToken *string,
) {
	if *poolContractAddress == "" {
		log.Fatal("missing required flag: pool-contract-address")
	}
	if *poolToken0Ticker == "" {
		log.Fatal("missing required flag: pool-token0-ticker")
	}
	if *poolToken0DecimalsNum == 0 {
		log.Fatal("missing required flag: pool-token0-decimals")
	}
	if *poolToken1Ticker == "" {
		log.Fatal("missing required flag: pool-token1-ticker")
	}
	if *poolToken1DecimalsNum == 0 {
		log.Fatal("missing required flag: pool-token1-decimals")
	}
	if *telegramChatID == 0 {
		log.Fatal("missing required flag: telegram-chat-id")
	}
	if *telegramBotToken == "" {
		log.Fatal("missing required flag: telegram-bot-token")
	}
}
