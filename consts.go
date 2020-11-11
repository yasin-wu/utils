package utils

type FileInfo struct {
	Name     string
	Path     string
	FileType string
	Content  string
}

const (
	Ali_SMS_Scheme  = "https"
	Ali_SMS_SUCCESS = "OK"
)

const (
	DayTimestamp int64 = 24 * 60 * 60 * 1000
)
