package test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	email2 "github.com/yasin-wu/utils/email"
)

func TestEmail(t *testing.T) {
	host := ""
	port := ""
	user := ""
	password := ""
	from := ""
	to := ""
	email, err := email2.New(host, port, user, password, from)
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
