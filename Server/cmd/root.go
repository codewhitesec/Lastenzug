package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

// nolint:gochecknoglobals
var rootCmd = &cobra.Command{
	Use:          "talje",
	SilenceUsage: true,
}

// nolint:gochecknoglobals
var mainContext context.Context

// Execute is the main cobra method
// copied from https://github.com/OJ/gobuster/blob/master/cli/cmd/root.go
func Execute() {
	var cancel context.CancelFunc
	mainContext, cancel = context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()
	go func() {
		select {
		case <-signalChan:
			// caught CTRL+C
			fmt.Println("\n[!] Keyboard interrupt detected, terminating.")
			cancel()
			os.Exit(1)
		case <-mainContext.Done():
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		// Leaving this in results in the same error appearing twice
		// Once before and once after the help output. Not sure if
		// this is going to be needed to output other errors that
		// aren't automatically outputted.
		// fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// specific flag parsing
	// log.SetFlags(log.Llongfile)
}
