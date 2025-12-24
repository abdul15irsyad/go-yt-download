package ffmpeg

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/abdul15irsyad/go-yt-download/pkg/models"
)

type FFMpeg struct {
	ffmpegPath string
}

func NewFFMpeg() *FFMpeg {
	return &FFMpeg{
		ffmpegPath: "ffmpeg", // ffmpeg in PATH
	}
}

func (f *FFMpeg) SetFFMpegPath(path string) {
	f.ffmpegPath = path
}

func (f *FFMpeg) Trim(req models.VideoTrimRequest) models.TrimResult {
	result := models.TrimResult{}

	// validation
	if req.InputPath == "" {
		result.Error = "input path cannot be empty"
		return result
	}
	if req.OutputPath == "" {
		result.Error = "output path cannot be empty"
		return result
	}
	if _, err := os.Stat(req.InputPath); os.IsNotExist(err) {
		result.Error = fmt.Sprintf("file input not found: %s", req.InputPath)
		return result
	}
	if err := validateTimeFormat(req.StartTime); err != nil {
		result.Error = fmt.Sprintf("start time invalid: %v", err)
		return result
	}
	if err := validateTimeFormat(req.EndTime); err != nil {
		result.Error = fmt.Sprintf("end time invalid: %v", err)
		return result
	}

	// build ffmpeg command
	fmt.Printf("trim video from %s to %s\n", req.StartTime, req.EndTime)
	cmd := exec.Command(
		f.ffmpegPath,
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

	// run command
	if err := cmd.Run(); err != nil {
		result.Error = fmt.Sprintf("trim video failed: %v", err)
		os.Remove(req.OutputPath)
		return result
	}

	result.Success = true
	result.OutputPath = req.OutputPath
	fmt.Printf("video trimmed: %s\n", req.OutputPath)

	return result
}

func (f *FFMpeg) TrimMultiple(req models.TrimSegment) []models.TrimResult {
	results := make([]models.TrimResult, len(req.Segments))

	// validation
	if req.InputPath == "" {
		results = append(results, models.TrimResult{
			Error: "input path cannot be empty",
		})
		return results
	}
	if _, err := os.Stat(req.InputPath); os.IsNotExist(err) {
		results = append(results, models.TrimResult{
			Error: fmt.Sprintf("input file not found: %s", req.InputPath),
		})
		return results
	}
	if len(req.Segments) == 0 {
		results = append(results, models.TrimResult{
			Error: "there are no segments to trim",
		})
		return results
	}

	sem := make(chan struct{}, 3)
	var wg sync.WaitGroup
	var rwMutex sync.RWMutex

	// process each segment
	for i, segment := range req.Segments {
		if segment.OutputPath == "" {
			results[i] = models.TrimResult{
				Error: fmt.Sprintf("output path for segment %d cannot be empty", i+1),
			}
			continue
		}
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer func() {
				<-sem
				wg.Done()
			}()
			fmt.Printf("processing segment %d of %d\n", i+1, len(req.Segments))

			trimReq := models.VideoTrimRequest{
				InputPath:  req.InputPath,
				OutputPath: segment.OutputPath,
				StartTime:  segment.StartTime,
				EndTime:    segment.EndTime,
			}

			result := f.Trim(trimReq)
			rwMutex.Lock()
			results[i] = result
			rwMutex.Unlock()
		}()
	}

	wg.Wait()

	return results
}

func (f *FFMpeg) MergeVideoAndAudio(videoFile, audioFile string) (io.Reader, *exec.Cmd, error) {
	cmd := exec.Command(
		f.ffmpegPath,
		"-i", videoFile,
		"-i", audioFile,
		"-c:v", "copy",
		"-c:a", "aac",
		"-movflags", "frag_keyframe+empty_moov",
		"-f", "mp4",
		"pipe:1",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	// cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	return stdout, cmd, nil
}

func validateTimeFormat(timeStr string) error {
	_, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return fmt.Errorf("time format invalid, use HH:MM:SS format")
	}
	return nil
}
