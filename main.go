package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/aqyuki/expand-bot/config"
	"github.com/aqyuki/expand-bot/discord"
)

type exitCode int

const (
	ExitCodeOK exitCode = iota
	ExitCodeError
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	code := run(ctx)
	defer exit(code)
}

func run(ctx context.Context) exitCode {
	store := config.NewStore()
	bot := discord.NewBot(store)

	if err := bot.Start(); err != nil {
		fmt.Printf("failed to start the bot\n\t%v\n", err)
		return ExitCodeError
	}

	<-ctx.Done()
	if err := bot.Stop(); err != nil {
		fmt.Printf("failed to stop the bot\n\t%v\n", err)
		return ExitCodeError
	}

	return ExitCodeOK
}

// exit is a wrapper of os.Exit.
func exit[T ~int](code T) {
	os.Exit(int(code))
}
