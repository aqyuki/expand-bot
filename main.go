package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/aqyuki/cham/discord"
	"github.com/aqyuki/cham/logging"
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
	logger := logging.DefaultLogger()
	bot := discord.NewBot(loadTokenFromEnv(), discord.WithLogger(logger))

	if err := bot.Start(); err != nil {
		logger.Errorf("failed to start the bot\n\t%v\n", err)
		return ExitCodeError
	}

	<-ctx.Done()
	fmt.Printf("received signal to stop the bot\nshutting down...\n")
	if err := bot.Stop(); err != nil {
		logger.Errorf("failed to stop the bot\n\t%v\n", err)
		return ExitCodeError
	}
	return ExitCodeOK
}

// exit is a wrapper of os.Exit.
func exit[T ~int](code T) {
	os.Exit(int(code))
}

func loadTokenFromEnv() string {
	return os.Getenv("DISCORD_TOKEN")
}
