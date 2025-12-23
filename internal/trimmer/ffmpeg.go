package trimmer

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/abdul15irsyad/go-yt-download/pkg/models"
)

// FFMpegTrimmer menangani pemotongan video menggunakan FFmpeg
type FFMpegTrimmer struct {
	ffmpegPath string
}

// NewFFMpegTrimmer membuat instance baru FFMpegTrimmer
func NewFFMpegTrimmer() *FFMpegTrimmer {
	return &FFMpegTrimmer{
		ffmpegPath: "ffmpeg", // Asumsi ffmpeg sudah di PATH
	}
}

// SetFFMpegPath mengatur path custom untuk ffmpeg
func (ft *FFMpegTrimmer) SetFFMpegPath(path string) {
	ft.ffmpegPath = path
}

// Trim memotong video berdasarkan start dan end time
func (ft *FFMpegTrimmer) Trim(req models.VideoTrimRequest) models.TrimResult {
	result := models.TrimResult{}

	// Validasi input
	if req.InputPath == "" {
		result.Error = "input path tidak boleh kosong"
		return result
	}

	if req.OutputPath == "" {
		result.Error = "output path tidak boleh kosong"
		return result
	}

	// Validasi file input
	if _, err := os.Stat(req.InputPath); os.IsNotExist(err) {
		result.Error = fmt.Sprintf("file input tidak ditemukan: %s", req.InputPath)
		return result
	}

	// Validasi waktu
	if err := validateTimeFormat(req.StartTime); err != nil {
		result.Error = fmt.Sprintf("format start time tidak valid: %v", err)
		return result
	}

	if err := validateTimeFormat(req.EndTime); err != nil {
		result.Error = fmt.Sprintf("format end time tidak valid: %v", err)
		return result
	}

	// Build FFmpeg command
	fmt.Printf("Memotong video dari %s ke %s\n", req.StartTime, req.EndTime)
	cmd := exec.Command(
		ft.ffmpegPath,
		"-i", req.InputPath,
		"-ss", req.StartTime,
		"-to", req.EndTime,
		"-c:v", "libx264",
		"-crf", "23",
		"-preset", "medium",
		"-c:a", "aac",
		"-b:a", "192k",
		"-y", // Overwrite output file
		req.OutputPath,
	)

	// Run command
	if err := cmd.Run(); err != nil {
		result.Error = fmt.Sprintf("gagal memotong video: %v", err)
		os.Remove(req.OutputPath) // Hapus file jika gagal
		return result
	}

	result.Success = true
	result.OutputPath = req.OutputPath
	fmt.Printf("✓ Video berhasil dipotong: %s\n", req.OutputPath)

	return result
}

// TrimMultiple memotong video menjadi multiple segments
func (ft *FFMpegTrimmer) TrimMultiple(req models.TrimSegment) []models.TrimResult {
	results := make([]models.TrimResult, 0)

	// Validasi input
	if req.InputPath == "" {
		results = append(results, models.TrimResult{
			Error: "input path tidak boleh kosong",
		})
		return results
	}

	if _, err := os.Stat(req.InputPath); os.IsNotExist(err) {
		results = append(results, models.TrimResult{
			Error: fmt.Sprintf("file input tidak ditemukan: %s", req.InputPath),
		})
		return results
	}

	if len(req.Segments) == 0 {
		results = append(results, models.TrimResult{
			Error: "tidak ada segment untuk dipotong",
		})
		return results
	}

	// Proses setiap segment
	for i, segment := range req.Segments {
		fmt.Printf("\nMemproses segment %d dari %d\n", i+1, len(req.Segments))

		if segment.OutputPath == "" {
			results = append(results, models.TrimResult{
				Error: fmt.Sprintf("output path untuk segment %d tidak boleh kosong", i+1),
			})
			continue
		}

		trimReq := models.VideoTrimRequest{
			InputPath:  req.InputPath,
			OutputPath: segment.OutputPath,
			StartTime:  segment.StartTime,
			EndTime:    segment.EndTime,
		}

		result := ft.Trim(trimReq)
		results = append(results, result)
	}

	return results
}

// GetVideoDuration mendapatkan durasi video
func (ft *FFMpegTrimmer) GetVideoDuration(videoPath string) (time.Duration, error) {
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		return 0, fmt.Errorf("file tidak ditemukan: %s", videoPath)
	}

	// FFprobe command untuk mendapatkan durasi
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1:nokey=1",
		videoPath,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("gagal mendapatkan durasi video: %v", err)
	}

	var seconds float64
	fmt.Sscanf(string(output), "%f", &seconds)

	return time.Duration(seconds) * time.Second, nil
}

// validateTimeFormat memvalidasi format waktu HH:MM:SS
func validateTimeFormat(timeStr string) error {
	_, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return fmt.Errorf("format waktu tidak valid, gunakan HH:MM:SS")
	}
	return nil
}
