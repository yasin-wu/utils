package common

import "log"

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
