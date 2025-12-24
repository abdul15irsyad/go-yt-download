package main

import (
	"fmt"

	"github.com/abdul15irsyad/go-yt-download/internal/downloader"
	"github.com/abdul15irsyad/go-yt-download/internal/ffmpeg"
	"github.com/abdul15irsyad/go-yt-download/pkg/models"
	"github.com/spf13/cobra"
)

var (
	downloadURL string
	downloadDir string
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download videos from YouTube",
	Long: `Download videos from YouTube with the given URL.
Example: yt-download download -u "https://www.youtube.com/watch?v=..."`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if downloadURL == "" {
			return fmt.Errorf("URL cannot be empty, use flag -u or --url")
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
	downloadCmd.Flags().StringVarP(&downloadDir, "output", "o", "./downloads", "output directory for video")
	downloadCmd.MarkFlagRequired("url")
	rootCmd.AddCommand(downloadCmd)
}
