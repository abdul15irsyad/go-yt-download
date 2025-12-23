# Go YouTube Downloader & Trimmer

Project Go untuk download YouTube video dan memotong video berdasarkan waktu tertentu.

## Fitur

- ✅ Download video dari YouTube
- ✅ Potong video dengan start dan end time tertentu
- ✅ Potong video menjadi multiple segments sekaligus
- ✅ CLI yang mudah digunakan dengan Cobra

## Prerequisites

1. **Go 1.21+** - [Download](https://golang.org/dl/)
2. **FFmpeg** - Untuk pemotongan video
   - Windows: Download dari [ffmpeg.org](https://ffmpeg.org/download.html)
   - Linux: `sudo apt-get install ffmpeg`
   - macOS: `brew install ffmpeg`
3. **FFprobe** - Untuk mendapatkan info video (biasanya sudah included dengan FFmpeg)

## Instalasi

### 1. Clone atau setup project

```bash
cd c:\Users\ABDUL\Documents\Code\go-yt-download
```

### 2. Download dependencies

```bash
go mod download
```

### 3. Build aplikasi

```bash
go build -o yt-download.exe ./cmd
```

## Penggunaan

### 1. Download Video dari YouTube

```bash
# Basic
yt-download download -u "https://www.youtube.com/watch?v=..."

# Dengan custom output directory
yt-download download -u "https://www.youtube.com/watch?v=..." -o ./my-videos
```

### 2. Potong Video

```bash
# Format: yt-download trim -i <input> -o <output> -s <start_time> -e <end_time>
yt-download trim -i video.mp4 -o output.mp4 -s 00:00:10 -e 00:00:30

# Contoh: Potong dari 10 detik sampai 30 detik
yt-download trim -i myVideo.mp4 -o trimmed.mp4 -s 00:00:10 -e 00:00:30

# Contoh: Potong dari 1 menit 30 detik sampai 3 menit
yt-download trim -i myVideo.mp4 -o trimmed.mp4 -s 00:01:30 -e 00:03:00
```

### 3. Potong Video menjadi Multiple Segments

Buat file `segments.json`:

```json
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
    },
    {
      "startTime": "00:01:00",
      "endTime": "00:01:30",
      "outputPath": "segment3.mp4"
    }
  ]
}
```

Jalankan:

```bash
yt-download trim-multiple -i video.mp4 -c segments.json
```

## Struktur Project

```txt
go-yt-download/
├── cmd/
│   ├── main.go              # Entry point
│   ├── download.go          # Command download
│   ├── trim.go              # Command trim single
│   └── trim_multiple.go     # Command trim multiple
├── internal/
│   ├── downloader/
│   │   └── youtube.go       # YouTube downloader logic
│   └── trimmer/
│       └── ffmpeg.go        # FFmpeg trimmer logic
├── pkg/
│   └── models/
│       └── models.go        # Data structures
├── go.mod
└── README.md
```

## API Documentation

### Models

#### VideoDownloadRequest

```go
type VideoDownloadRequest struct {
 URL       string  // URL video YouTube
 OutputDir string  // Directory output
}
```

#### VideoTrimRequest

```go
type VideoTrimRequest struct {
 InputPath  string  // Path video input
 OutputPath string  // Path video output
 StartTime  string  // HH:MM:SS
 EndTime    string  // HH:MM:SS
}
```

#### TrimSegment

```go
type TrimSegment struct {
 InputPath string      // Path video input
 Segments  []Segment   // List segments
}

type Segment struct {
 StartTime  string  // HH:MM:SS
 EndTime    string  // HH:MM:SS
 OutputPath string  // Path output
}
```

### Functions

#### YouTubeDownloader

- `NewYouTubeDownloader() *YouTubeDownloader`
- `Download(req VideoDownloadRequest) DownloadResult`

#### FFMpegTrimmer

- `NewFFMpegTrimmer() *FFMpegTrimmer`
- `SetFFMpegPath(path string)`
- `Trim(req VideoTrimRequest) TrimResult`
- `TrimMultiple(req TrimSegment) []TrimResult`
- `GetVideoDuration(videoPath string) (time.Duration, error)`

## Contoh Program Programmatic

```go
package main

import (
 "fmt"
 "github.com/abdul15irsyad/go-yt-download/internal/downloader"
 "github.com/abdul15irsyad/go-yt-download/internal/trimmer"
 "github.com/abdul15irsyad/go-yt-download/pkg/models"
)

func main() {
 // Download video
 yd := downloader.NewYouTubeDownloader()
 dlResult := yd.Download(models.VideoDownloadRequest{
  URL:       "https://www.youtube.com/watch?v=...",
  OutputDir: "./videos",
 })

 if !dlResult.Success {
  fmt.Printf("Error: %s\n", dlResult.Error)
  return
 }

 // Potong video
 ft := trimmer.NewFFMpegTrimmer()
 trimResult := ft.Trim(models.VideoTrimRequest{
  InputPath:  dlResult.VideoPath,
  OutputPath: "trimmed.mp4",
  StartTime:  "00:00:10",
  EndTime:    "00:00:30",
 })

 if trimResult.Success {
  fmt.Printf("✓ Video berhasil dipotong: %s\n", trimResult.OutputPath)
 } else {
  fmt.Printf("Error: %s\n", trimResult.Error)
 }
}
```

## Troubleshooting

### FFmpeg tidak ditemukan

- Pastikan FFmpeg sudah diinstall dan ada di PATH
- Test dengan: `ffmpeg -version`

### YouTube download gagal

- Periksa koneksi internet
- Verifikasi URL video valid dan bisa diakses
- Beberapa video mungkin memiliki restriction

### Video trim gagal

- Pastikan format waktu benar (HH:MM:SS)
- Verifikasi end time > start time
- Pastikan file input ada dan bisa dibaca

## Development

### Run without building

```bash
go run ./cmd download -u "https://www.youtube.com/watch?v=..."
go run ./cmd trim -i input.mp4 -o output.mp4 -s 00:00:10 -e 00:00:30
```

### Test

```bash
go test ./...
```

## License

MIT

## Contributing

Feel free to submit issues dan pull requests!
