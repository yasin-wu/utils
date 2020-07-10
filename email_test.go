package utils

import (
	"testing"
)

func TestEmail_Send(t *testing.T) {
	smtpServer := &SMTPServer{
		Host:     "smtp.qq.com",
		Port:     "25",
		User:     "yasin_wu@qq.com",
		Password: "mjfvvjhqmrjocajd",
	}

	email := &Email{
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
