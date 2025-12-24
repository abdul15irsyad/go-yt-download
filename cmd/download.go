package main

import (
	"fmt"
	"log"
	"os"

	"github.com/abdul15irsyad/go-yt-download/internal/downloader"
	"github.com/abdul15irsyad/go-yt-download/internal/ffmpeg"
	"github.com/abdul15irsyad/go-yt-download/pkg/models"
	"github.com/spf13/cobra"
)

var (
	downloadURL             string
	downloadDir             string
	defaultOuputDownloadDir = "./output/downloads"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download videos from YouTube",
	Long: `Download videos from YouTube with the given URL.
Example: yt-download download -u "https://www.youtube.com/watch?v=..."`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := os.MkdirAll(defaultOuputDownloadDir, 0755); err != nil {
			log.Fatal("failed to create default output directory: %v", err)
		}
		if downloadURL == "" {
			return fmt.Errorf("url cannot be empty, use flag -u or --url")
		}
		if downloadDir == "" {
			downloadDir = defaultOuputDownloadDir
		}
		ffmpeg := ffmpeg.NewFFMpeg()
		yd := downloader.NewYouTubeDownloader(ffmpeg)
		req := models.VideoDownloadRequest{
			URL:       downloadURL,
			OutputDir: downloadDir,
		}

		result := yd.Download(req)
		if !result.Success {
			return fmt.Errorf("download failed: %s", result.Error)
		}

		return nil
	},
}

func init() {
	downloadCmd.Flags().StringVarP(&downloadURL, "url", "u", "", "URL video youTube (required)")
	downloadCmd.Flags().StringVarP(&downloadDir, "output", "o", "", "output directory for video")
	downloadCmd.MarkFlagRequired("url")
	rootCmd.AddCommand(downloadCmd)
}
