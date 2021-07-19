package test

import (
	"testing"

	email2 "github.com/yasin-wu/utils/email"
)

func TestEmail_Send(t *testing.T) {
	smtpServer := &email2.SMTPServer{
		Host:     "smtp.qq.com",
		Port:     "25",
		User:     "yasin_wu@qq.com",
		Password: "mjfvvjhqmrjocajd",
	}

	email := &email2.Email{
		From:       "yasin_wu@qq.com",
		To:         "yasin_wu@qq.com",
		Subject:    "test",
		Content:    "test",
		SMTPServer: smtpServer,
	}

	err := email.Send()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("send email ok")
}
