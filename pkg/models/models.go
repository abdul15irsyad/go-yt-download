package models

type VideoDownloadRequest struct {
	URL       string
	OutputDir string
}

type VideoTrimRequest struct {
	InputPath  string
	OutputPath string
	StartTime  string // format: HH:MM:SS
	EndTime    string // format: HH:MM:SS
}

type TrimSegment struct {
	InputPath string
	Segments  []Segment
}

type Segment struct {
	StartTime  string // format: HH:MM:SS
	EndTime    string // format: HH:MM:SS
	OutputPath string
}

type DownloadResult struct {
	Success   bool
	VideoPath string
	Error     string
}

type TrimResult struct {
	Success    bool
	OutputPath string
	Error      string
}
