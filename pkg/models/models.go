package models

// VideoDownloadRequest berisi informasi untuk download video
type VideoDownloadRequest struct {
	URL       string
	OutputDir string
}

// VideoTrimRequest berisi informasi untuk memotong video
type VideoTrimRequest struct {
	InputPath  string
	OutputPath string
	StartTime  string // format: HH:MM:SS
	EndTime    string // format: HH:MM:SS
}

// TrimSegment berisi multiple segment yang akan dipotong
type TrimSegment struct {
	InputPath string
	Segments  []Segment
}

// Segment berisi start dan end time untuk satu potongan
type Segment struct {
	StartTime  string // format: HH:MM:SS
	EndTime    string // format: HH:MM:SS
	OutputPath string
}

// DownloadResult berisi hasil download
type DownloadResult struct {
	Success   bool
	VideoPath string
	Error     string
}

// TrimResult berisi hasil trimming
type TrimResult struct {
	Success    bool
	OutputPath string
	Error      string
}
