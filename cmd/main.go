package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yt-download",
	Short: "Download YouTube videos and trim them based on specific time",
	Long: `yt-download CLI application for:
1. Download videos from YouTube
2. Trim video based on specific timestamps`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
