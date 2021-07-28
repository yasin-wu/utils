package test

import (
	"testing"

	email2 "github.com/yasin-wu/utils/email"
)

func TestEmailSend(t *testing.T) {
	email, err := email2.New("smtp.qq.com", "25",
		"yasin_wu@qq.com", "gumrjpxqvnqrbhai", "yasin_wu@qq.com")
	if err != nil {
		t.Error(err)
		return
	}
	err = email.Send("yasin_wu@qq.com", "test", "test")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("send email ok")
}
