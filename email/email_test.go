package email

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestEmail(t *testing.T) {
	to := ""
	email, err := New(&Config{})
	if err != nil {
		log.Fatal(err)
	}
	tos := strings.Split(to, ",")
	err = email.SendTLS(tos, "test", "test")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("send email ok")
}
