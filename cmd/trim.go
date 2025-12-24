package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/abdul15irsyad/go-yt-download/internal/ffmpeg"
	"github.com/abdul15irsyad/go-yt-download/pkg/models"
	"github.com/spf13/cobra"
)

var (
	trimInput           string
	trimOutput          string
	trimTimeRanges      string
	defaultOuputTrimDir = "./output/trims"
)

var trimCmd = &cobra.Command{
	Use:   "trim",
	Short: "Trim videos based on specific time ranges",
	Long: `Trim videos by specifying one or more time ranges.
time format: HH:MM:SS (example: 00:00:10 for 10 seconds)
multiple ranges are separated by commas.
example: yt-download trim -i input.mp4 -o output.mp4 -t 00:00:05-00:00:30
example with multiple ranges: yt-download trim -i input.mp4 -o output.mp4 -t 00:00:05-00:00:30,00:01:00-00:01:30`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := os.MkdirAll(defaultOuputTrimDir, 0755); err != nil {
			log.Fatal("failed to create default output directory: %v", err)
		}

		// validation
		if trimInput == "" {
			return fmt.Errorf("the input file cannot be empty, use the -i or --input flag")
		}
		if trimOutput == "" {
			fileName := filepath.Base(trimInput)
			trimOutput = fmt.Sprintf("%s/%d_%s", defaultOuputTrimDir, time.Now().UnixMilli(), fileName)
			fmt.Println(trimOutput)
		}
		if trimTimeRanges == "" {
			return fmt.Errorf("time ranges cannot be empty, use the -t or --time flag")
		}

		// parse time ranges
		segments, err := parseTimeRanges(trimTimeRanges, trimOutput)
		if err != nil {
			return err
		}

		// ff single segment, use simple trim
		if len(segments) == 1 {
			ft := ffmpeg.NewFFMpeg()
			req := models.VideoTrimRequest{
				InputPath:  trimInput,
				OutputPath: trimOutput,
				StartTime:  segments[0].StartTime,
				EndTime:    segments[0].EndTime,
			}

			result := ft.Trim(req)
			if !result.Success {
				return fmt.Errorf("trim failed: %s", result.Error)
			}

			return nil
		}

		// for multiple segments, use TrimMultiple
		ft := ffmpeg.NewFFMpeg()
		req := models.TrimSegment{
			InputPath: trimInput,
			Segments:  segments,
		}

		results := ft.TrimMultiple(req)

		// print results
		fmt.Println("\n=== results ===")
		successCount := 0
		for i, result := range results {
			if result.Success {
				fmt.Printf("segment %d: success - %s\n", i+1, result.OutputPath)
				successCount++
			} else {
				fmt.Printf("segment %d: failed - %s\n", i+1, result.Error)
			}
		}
		fmt.Printf("\ntotal: %d/%d successful\n", successCount, len(results))

		if successCount < len(results) {
			return fmt.Errorf("some segments failed to process")
		}

		return nil
	},
}

// parses time ranges in format: start-end,start-end,...
func parseTimeRanges(timeRangesStr string, outputPath string) ([]models.Segment, error) {
	var segments []models.Segment
	ranges := strings.Split(timeRangesStr, ",")

	for i, timeRange := range ranges {
		timeRange = strings.TrimSpace(timeRange)
		parts := strings.Split(timeRange, "-")

		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid time range format: %s (expected format: HH:MM:SS-HH:MM:SS)", timeRange)
		}

		startTime := strings.TrimSpace(parts[0])
		endTime := strings.TrimSpace(parts[1])
		if err := validateTimeFormat(startTime); err != nil {
			return nil, fmt.Errorf("invalid start time format: %s - %v", startTime, err)
		}
		if err := validateTimeFormat(endTime); err != nil {
			return nil, fmt.Errorf("invalid end time format: %s - %v", endTime, err)
		}

		// for single segment, use the provided output path
		// for multiple segments, append segment number to output path
		var outPath string
		if len(ranges) == 1 {
			outPath = outputPath
		} else {
			ext := ""
			basePath := outputPath
			if lastDot := strings.LastIndex(outputPath, "."); lastDot != -1 {
				ext = outputPath[lastDot:]
				basePath = outputPath[:lastDot]
			}
			outPath = fmt.Sprintf("%s_%d%s", basePath, i+1, ext)
		}

		segments = append(segments, models.Segment{
			StartTime:  startTime,
			EndTime:    endTime,
			OutputPath: outPath,
		})
	}

	return segments, nil
}

// validates HH:MM:SS format
func validateTimeFormat(timeStr string) error {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return fmt.Errorf("expected HH:MM:SS format")
	}
	return nil
}

func init() {
	trimCmd.Flags().StringVarP(&trimInput, "input", "i", "", "Path to video input file (required)")
	trimCmd.Flags().StringVarP(&trimOutput, "output", "o", "", "Path to video output file (required)")
	trimCmd.Flags().StringVarP(&trimTimeRanges, "time", "t", "", "Time ranges in format: HH:MM:SS-HH:MM:SS[,HH:MM:SS-HH:MM:SS,...] (required)")

	trimCmd.MarkFlagRequired("input")
	trimCmd.MarkFlagRequired("time")

	rootCmd.AddCommand(trimCmd)
}
