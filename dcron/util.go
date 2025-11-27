package dcron

import (
	"fmt"
	"os"
	"time"
)

func getInstanceID() string {
	if hostname, err := os.Hostname(); err == nil {
		return hostname
	}
	return fmt.Sprintf("instance-%d", time.Now().UnixNano())
}
