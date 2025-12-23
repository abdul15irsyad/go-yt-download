package main

import (
	"fmt"

	"github.com/abdul15irsyad/go-yt-download/internal/trimmer"
	"github.com/abdul15irsyad/go-yt-download/pkg/models"
	"github.com/spf13/cobra"
)

var (
	trimInput  string
	trimOutput string
	trimStart  string
	trimEnd    string
)

var trimCmd = &cobra.Command{
	Use:   "trim",
	Short: "Potong video berdasarkan waktu tertentu",
	Long: `Potong video dengan menentukan start time dan end time.
Format waktu: HH:MM:SS (contoh: 00:00:10 untuk 10 detik)
Contoh: yt-download trim -i input.mp4 -o output.mp4 -s 00:00:05 -e 00:00:30`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validasi input
		if trimInput == "" {
			return fmt.Errorf("input file tidak boleh kosong, gunakan flag -i atau --input")
		}
		if trimOutput == "" {
			return fmt.Errorf("output file tidak boleh kosong, gunakan flag -o atau --output")
		}
		if trimStart == "" {
			return fmt.Errorf("start time tidak boleh kosong, gunakan flag -s atau --start")
		}
		if trimEnd == "" {
			return fmt.Errorf("end time tidak boleh kosong, gunakan flag -e atau --end")
		}

		ft := trimmer.NewFFMpegTrimmer()
		req := models.VideoTrimRequest{
			InputPath:  trimInput,
			OutputPath: trimOutput,
			StartTime:  trimStart,
			EndTime:    trimEnd,
		}

		result := ft.Trim(req)
		if !result.Success {
			return fmt.Errorf("trim gagal: %s", result.Error)
		}

		return nil
	},
}

func init() {
	trimCmd.Flags().StringVarP(&trimInput, "input", "i", "", "Path file video input (required)")
	trimCmd.Flags().StringVarP(&trimOutput, "output", "o", "", "Path file video output (required)")
	trimCmd.Flags().StringVarP(&trimStart, "start", "s", "", "Start time dalam format HH:MM:SS (required)")
	trimCmd.Flags().StringVarP(&trimEnd, "end", "e", "", "End time dalam format HH:MM:SS (required)")

	trimCmd.MarkFlagRequired("input")
	trimCmd.MarkFlagRequired("output")
	trimCmd.MarkFlagRequired("start")
	trimCmd.MarkFlagRequired("end")

	rootCmd.AddCommand(trimCmd)
}
