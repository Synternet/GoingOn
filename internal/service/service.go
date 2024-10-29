package service

import (
	"context"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"syscall"

	svcn "github.com/synternet/goingon/pkg/nats"
	types "github.com/synternet/goingon/pkg/types"

	"github.com/synternet/goingon/internal/ethereum"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/sync/errgroup"
)

//go:embed abi.json
var abiDefinition string

type Token struct {
	Ticker   string
	Decimals *big.Int
}

type Config struct {
	TelegramBotToken string
	TelegramChatID   int64
	ContractAddress  string
	Token0           Token
	Token1           Token
}

type Service struct {
	abi         abi.ABI
	ctx         context.Context
	cfg         Config
	nats        *svcn.NatsService
	telegramBot *tgbotapi.BotAPI
}

func NewService(s *svcn.NatsService, ctx context.Context, cfg Config) *Service {
	ABI, err := abi.JSON(strings.NewReader(abiDefinition))
	if err != nil {
		log.Fatalf("failed to parse ABI: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatalf("failed to create telegram bot: %v", err)
	}

	return &Service{
		ctx:         ctx,
		cfg:         cfg,
		nats:        s,
		abi:         ABI,
		telegramBot: bot,
	}
}

func (s Service) ProcessTxLogEventFromStream(data []byte) error {
	incoming := types.EthLogEvent{}
	err := json.Unmarshal(data, &incoming)
	if err != nil {
		return err
	}

	if incoming.Address != s.cfg.ContractAddress {
		return nil
	}

	eventTicker, err := ethereum.GetEventName(incoming)
	if err != nil {
		// Not an exhaustive events list. Silently ignore unknown.
		return nil
	}

	eventData, err := hex.DecodeString(strings.TrimPrefix(incoming.Data, "0x"))
	if err != nil {
		log.Fatalf("failed to decode log data: %v", err)
		return nil
	}

	switch eventTicker {
	case "Swap":
		fmt.Println(incoming)
		swap := new(types.Swap)
		err = s.abi.UnpackIntoInterface(swap, "Swap", eventData)
		if err != nil {
			log.Fatalf("failed to decode Swap event log: %v", err)
		}

		var msg string
		if swap.Amount0In.Sign() == 0 {
			msg = s.formatSwapMessage(incoming.Topics[2], swap.Amount1In, swap.Amount0Out, s.cfg.Token1.Decimals, s.cfg.Token0.Decimals, s.cfg.Token1.Ticker, s.cfg.Token0.Ticker, incoming.TransactionHash, s.cfg.ContractAddress, s.cfg.Token0.Ticker, s.cfg.Token1.Ticker, true)
		} else {
			msg = s.formatSwapMessage(incoming.Topics[2], swap.Amount0In, swap.Amount1Out, s.cfg.Token0.Decimals, s.cfg.Token1.Decimals, s.cfg.Token0.Ticker, s.cfg.Token1.Ticker, incoming.TransactionHash, s.cfg.ContractAddress, s.cfg.Token0.Ticker, s.cfg.Token1.Ticker, false)
		}
		s.sendMessageToTelegram(msg)
		log.Print(msg)
	}

	return nil
}

func (s Service) Serve() {
	serveCtx, cancelFn := context.WithCancel(s.ctx)
	defer cancelFn()

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Service is interrupted")
		cancelFn()
	}()

	rungroup, groupCtx := errgroup.WithContext(serveCtx)

	if s.nats != nil {
		rungroup.Go(func() error {
			return s.nats.Serve(groupCtx)
		})
	}

	log.Println("Service is started")

	if err := rungroup.Wait(); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("Service is stopped %s", err.Error())
		}
	}

	var completionGroup errgroup.Group
	if s.nats != nil {
		completionGroup.Go(func() error {
			return nil
		})
	}
}

func (s *Service) sendMessageToTelegram(text string) {
	msg := tgbotapi.NewMessage(s.cfg.TelegramChatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.DisableWebPagePreview = true
	_, err := s.telegramBot.Send(msg)
	if err != nil {
		log.Printf("failed to send telegram message: %v", err)
	}
}

func (s *Service) formatSwapMessage(to string, amountIn *big.Int, amountOut *big.Int, amountInDecimals *big.Int, amountOutDecimals *big.Int, tickerInTicker string, tickerOutTicker string, txHash string, contractAddress string, ticker0 string, ticker1 string, buy bool) string {
	fTo := FormatAddress(common.BytesToAddress(common.HexToHash(to).Bytes()).String())
	fAmountIn := new(big.Float).Quo(new(big.Float).SetInt(amountIn), new(big.Float).SetInt(amountInDecimals))
	fAmountOut := new(big.Float).Quo(new(big.Float).SetInt(amountOut), new(big.Float).SetInt(amountOutDecimals))
	link := fmt.Sprintf("https://etherscan.io/tx/%s", txHash)
	linkTo := fmt.Sprintf("https://etherscan.io/address/%s", common.BytesToAddress(common.HexToHash(to).Bytes()).String())
	poolPair := fmt.Sprintf("%s-%s", ticker0, ticker1)
	fTxHash := FormatAddress(txHash)
	tickerIn := tickerInTicker
	tickerOut := tickerOutTicker
	var icon string
	if buy {
		icon = fmt.Sprintf("🟢 Bought %s", ticker0)
	} else {
		icon = fmt.Sprintf("🔴 Sold %s", ticker0)
	}
	fStr := `Account: [%s](%s)
%s
%s %f ➡ %s %f
On [Uniswap %s](https://v2.info.uniswap.org/pair/%s) @ [%s](%s)`
	return fmt.Sprintf(fStr, fTo, linkTo, icon, tickerIn, fAmountIn, tickerOut, fAmountOut, poolPair, contractAddress, fTxHash, link)
}
