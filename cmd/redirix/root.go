package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulnebify/redirix/internal/app"
)

const asciiArt = `
 ____   _____  ____   ___  ____   ___ __  __
|  _ \ | ____||  _ \ |_ _||  _ \ |_ _|\ \/ /
| |_) ||  _|  | | | | | | | |_) | | |  \  / 
|  _ < | |___ | |_| | | | |  _ <  | |  /  \ 
|_| \_\|_____||____/ |___||_| \_\|___|/_/\_\
										   
`

var rootCmd = &cobra.Command{
	Use:     "Redirix",
	Short:   "Redirix is a SOCKS5 proxy server that registers in Redis for dynamic use.",
	Long:    "Redirix is a SOCKS5 proxy server that registers in Redis for dynamic use in distributed systems. See more https://github.com/vulnebify/redirix",
	Version: app.Version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(asciiArt)
		_ = cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.AddCommand(serveCmd)
}
