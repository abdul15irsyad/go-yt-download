package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/abdul15irsyad/go-yt-download/internal/trimmer"
	"github.com/abdul15irsyad/go-yt-download/pkg/models"
	"github.com/spf13/cobra"
)

var (
	trimMultipleInput  string
	trimMultipleConfig string
)

var trimMultipleCmd = &cobra.Command{
	Use:   "trim-multiple",
	Short: "Potong video menjadi multiple segments",
	Long: `Potong video menjadi multiple segments sekaligus dengan config JSON.
Contoh config.json:
{
  "segments": [
    {
      "startTime": "00:00:00",
      "endTime": "00:00:30",
      "outputPath": "segment1.mp4"
    },
    {
      "startTime": "00:00:30",
      "endTime": "00:01:00",
      "outputPath": "segment2.mp4"
    }
  ]
}

Contoh: yt-download trim-multiple -i input.mp4 -c config.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validasi input
		if trimMultipleInput == "" {
			return fmt.Errorf("input file tidak boleh kosong, gunakan flag -i atau --input")
		}
		if trimMultipleConfig == "" {
			return fmt.Errorf("config file tidak boleh kosong, gunakan flag -c atau --config")
		}

		// Baca config file
		configData, err := os.ReadFile(trimMultipleConfig)
		if err != nil {
			return fmt.Errorf("gagal membaca config file: %v", err)
		}

		// Parse config JSON
		var configSegments struct {
			Segments []struct {
				StartTime  string `json:"startTime"`
				EndTime    string `json:"endTime"`
				OutputPath string `json:"outputPath"`
			} `json:"segments"`
		}

		if err := json.Unmarshal(configData, &configSegments); err != nil {
			return fmt.Errorf("gagal parse config JSON: %v", err)
		}

		// Convert ke models.Segment
		segments := make([]models.Segment, len(configSegments.Segments))
		for i, seg := range configSegments.Segments {
			segments[i] = models.Segment{
				StartTime:  seg.StartTime,
				EndTime:    seg.EndTime,
				OutputPath: seg.OutputPath,
			}
		}

		// Trim multiple
		ft := trimmer.NewFFMpegTrimmer()
		req := models.TrimSegment{
			InputPath: trimMultipleInput,
			Segments:  segments,
		}

		results := ft.TrimMultiple(req)

		// Print results
		fmt.Println("\n=== HASIL PEMOTONGAN ===")
		successCount := 0
		for i, result := range results {
			if result.Success {
				fmt.Printf("✓ Segment %d: Berhasil - %s\n", i+1, result.OutputPath)
				successCount++
			} else {
				fmt.Printf("✗ Segment %d: Gagal - %s\n", i+1, result.Error)
			}
		}
		fmt.Printf("\nTotal: %d/%d berhasil\n", successCount, len(results))

		if successCount < len(results) {
			return fmt.Errorf("beberapa segment gagal diproses")
		}

		return nil
	},
}

func init() {
	trimMultipleCmd.Flags().StringVarP(&trimMultipleInput, "input", "i", "", "Path file video input (required)")
	trimMultipleCmd.Flags().StringVarP(&trimMultipleConfig, "config", "c", "", "Path file config JSON (required)")

	trimMultipleCmd.MarkFlagRequired("input")
	trimMultipleCmd.MarkFlagRequired("config")

	rootCmd.AddCommand(trimMultipleCmd)
}
