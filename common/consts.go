package common

type FileInfo struct {
	Name     string
	Path     string
	FileType string
	Content  string
}

const (
	DayTimestamp int64 = 24 * 60 * 60 * 1000
)

const (
	Byte = 1 << (10 * iota)
	KB
	MB
	GB
	TB
	PB
)

const (
	WordRatio = 0.05
)

const (
	TopLeft int = iota
	TopRight
	BottomLeft
	BottomRight
	Center
)
