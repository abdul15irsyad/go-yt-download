package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yt-download",
	Short: "Download YouTube video dan potong berdasarkan waktu tertentu",
	Long: `yt-download adalah aplikasi CLI untuk:
1. Download video dari YouTube
2. Memotong video berdasarkan timestamp tertentu
3. Memotong video menjadi multiple segments`,
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
