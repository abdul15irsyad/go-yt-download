package downloader

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/abdul15irsyad/go-yt-download/internal/ffmpeg"
	"github.com/abdul15irsyad/go-yt-download/pkg/models"
	"github.com/kkdai/youtube/v2"
)

type YouTubeDownloader struct {
	client *youtube.Client
	ffmpeg *ffmpeg.FFMpeg
}

func NewYouTubeDownloader(ffmpeg *ffmpeg.FFMpeg) *YouTubeDownloader {
	return &YouTubeDownloader{
		client: &youtube.Client{},
		ffmpeg: ffmpeg,
	}
}

func (yd *YouTubeDownloader) Download(req models.VideoDownloadRequest) models.DownloadResult {
	result := models.DownloadResult{}

	// validation
	if req.URL == "" {
		result.Error = "url cannot be empty"
		return result
	}
	if req.OutputDir == "" {
		req.OutputDir = "./downloads"
	}
	if err := os.MkdirAll(req.OutputDir, 0755); err != nil {
		result.Error = fmt.Sprintf("failed to create directory: %v", err)
		return result
	}

	client := yd.client

	// get video info
	fmt.Printf("get video information from: '%s'\n", req.URL)
	video, err := yd.client.GetVideo(req.URL)
	if err != nil {
		result.Error = fmt.Sprintf("failed to retrieve video info: %v", err)
		return result
	}

	formats := video.Formats
	if len(formats) == 0 {
		result.Error = "no video formats available"
		return result
	}
	videoFormat := formats[0]

	audioFormats := video.Formats.WithAudioChannels()
	if len(audioFormats) == 0 {
		log.Fatal("audio format not found")
	}
	audioFormat := audioFormats[0]

	unixTimestamp := time.Now().UnixMilli()
	videoFilePath := filepath.Join(os.TempDir(), fmt.Sprintf("%d_video.mp4", unixTimestamp))
	audioFilePath := filepath.Join(os.TempDir(), fmt.Sprintf("%d_audio.m4a", unixTimestamp))

	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Println("downloading video...")
		if err := download(client, video, &videoFormat, videoFilePath); err != nil {
			errCh <- err
		}
		fmt.Println("video downloaded")
	}()

	go func() {
		defer wg.Done()
		fmt.Println("downloading audio...")
		if err := download(client, video, &audioFormat, audioFilePath); err != nil {
			errCh <- err
		}
		fmt.Println("audio downloaded")
	}()

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("start merging video & audio\n")
	reader, cmd, err := yd.ffmpeg.MergeVideoAndAudio(videoFilePath, audioFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = cmd.Wait()
		os.Remove(videoFilePath)
		os.Remove(audioFilePath)
	}()

	// create output file
	filename := sanitizeFilename(video.Title) + ".mp4"
	outputPath := filepath.Join(req.OutputDir, filename)

	file, err := os.Create(outputPath)
	if err != nil {
		result.Error = fmt.Sprintf("failed create output file: %v", err)
		return result
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		result.Error = fmt.Sprintf("failed save video: %v", err)
		os.Remove(outputPath)
		return result
	}

	result.Success = true
	result.VideoPath = outputPath
	fmt.Printf("video downloaded: %s\n", outputPath)

	os.Remove(videoFilePath)
	os.Remove(audioFilePath)

	return result
}

func download(client *youtube.Client, video *youtube.Video, format *youtube.Format, filename string) error {
	stream, _, err := client.GetStream(video, format)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return err
	}

	return nil
}

func sanitizeFilename(filename string) string {
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := filename

	for _, char := range invalidChars {
		result = strings.ReplaceAll(result, char, "_")
	}

	if len(result) > 200 {
		result = result[:200]
	}

	return result
}
