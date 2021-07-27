package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"regexp"
	"strings"

	"github.com/jordan-wright/email"
)

type Email struct {
	Host     string
	Port     string
	User     string
	PassWord string
	From     string
}

var (
	emailRegexpStr = `^[_a-z0-9-]+(\.[_a-z0-9-]+)*@[a-z0-9-]+(\.[a-z0-9-]+)*(\.[a-z]{2,})$`
)

func New(host, port, user, password, from string) (*Email, error) {
	if host == "" {
		return nil, errors.New("smtp server host is nil")
	}
	if port == "" {
		return nil, errors.New("smtp server port is nil")
	}
	if user == "" {
		return nil, errors.New("smtp server user is nil")
	}
	if password == "" {
		return nil, errors.New("smtp server password is nil")
	}
	return &Email{Host: host, Port: port, User: user, PassWord: password, From: from}, nil
}

func (this *Email) Send(to, subject, content string) error {
	err := this.check(to, subject, content)
	if err != nil {
		return err
	}
	return this.sendMail(to, subject, content)
}

func (this *Email) SendTLS(to, subject, content string) error {
	err := this.check(to, subject, content)
	if err != nil {
		return err
	}
	return this.sendTLSMail(to, subject, content)
}

func (this *Email) sendMail(to, subject, content string) error {
	e := email.NewEmail()
	e.From = this.From
	e.To = []string{to}
	e.Subject = subject
	e.Text = []byte(content)
	err := e.Send(fmt.Sprintf("%s:%s", this.Host, this.Port),
		smtp.PlainAuth("", this.User, this.PassWord, this.Host))
	if err != nil {
		return err
	}
	return nil
}

func (this *Email) sendTLSMail(to, subject, content string) error {
	header := make(map[string]string)
	header["From"] = this.From
	header["To"] = to
	header["Subject"] = subject
	header["Content-Type"] = "text/html; charset=UTF-8"
	body := content
	sendMsg := ""
	for k, v := range header {
		sendMsg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	sendMsg += "\r\n" + body
	err := this.sendMailUsingTLS(
		fmt.Sprintf("%s:%s", this.Host, this.Port),
		smtp.PlainAuth("", this.User, this.PassWord, this.Host),
		this.User,
		[]string{to},
		[]byte(sendMsg),
	)
	if err != nil {
		return err
	}
	return nil
}

func (this *Email) sendMailUsingTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) (err error) {
	c, err := this.dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func (this *Email) dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

func (this *Email) check(to, subject, content string) error {
	if to == "" {
		return errors.New("to email is nil")
	}
	if subject == "" {
		return errors.New("subject of email is nil")
	}
	if content == "" {
		return errors.New("content of email is nil")
	}
	if !this.isEmail(to) {
		return errors.New("to email is error")
	}
	return nil
}

func (this *Email) isEmail(to string) bool {
	if !strings.Contains(to, "@") {
		return false
	}
	emailRegexp, _ := regexp.Compile(emailRegexpStr)
	return emailRegexp.MatchString(to)
}
