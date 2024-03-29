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

const emailRegexpStr = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`

type Email struct {
	config *Config
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
}

func New(config *Config) (*Email, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}
	if config.Host == "" {
		return nil, errors.New("smtp server host is nil")
	}
	if config.Port == "" {
		return nil, errors.New("smtp server port is nil")
	}
	if config.User == "" {
		return nil, errors.New("smtp server user is nil")
	}
	if config.Password == "" {
		return nil, errors.New("smtp server password is nil")
	}
	return &Email{config: config}, nil
}

func (e *Email) Send(to []string, subject, content string) error {
	err := e.check(to, subject, content)
	if err != nil {
		return err
	}
	return e.sendMail(to, subject, content)
}

func (e *Email) SendTLS(to []string, subject, content string) error {
	err := e.check(to, subject, content)
	if err != nil {
		return err
	}
	return e.sendTLSMail(to, subject, content)
}

func (e *Email) sendMail(to []string, subject, content string) error {
	el := email.NewEmail()
	el.From = e.config.From
	el.To = to
	el.Subject = subject
	el.Text = []byte(content)
	err := el.Send(e.addr(), e.plainAuth())
	if err != nil {
		return err
	}
	return nil
}

func (e *Email) sendTLSMail(to []string, subject, content string) error {
	header := make(map[string]any)
	header["From"] = e.config.From
	header["Subject"] = subject
	header["Content-Type"] = "text/html; charset=UTF-8"
	body := content
	sendMsg := ""
	for k, v := range header {
		sendMsg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	sendMsg += "\r\n" + body
	err := e.sendMailUsingTLS(to, []byte(sendMsg))
	if err != nil {
		return err
	}
	return nil
}

func (e *Email) sendMailUsingTLS(to []string, msg []byte) (err error) {
	addr := e.addr()
	auth := e.plainAuth()
	client, err := e.dial(addr)
	if err != nil {
		return err
	}
	defer e.cloesClient(client)
	if auth != nil {
		if ok, _ := client.Extension("AUTH"); ok {
			if err = client.Auth(auth); err != nil {
				return err
			}
		}
	}
	if err = client.Mail(e.config.User); err != nil {
		return err
	}
	for _, v := range to {
		if err = client.Rcpt(v); err != nil {
			return err
		}
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return client.Quit()
}

func (e *Email) dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

func (e *Email) check(to []string, subject, content string) error {
	if to == nil {
		return errors.New("to email is nil")
	}
	if subject == "" {
		return errors.New("subject of email is nil")
	}
	if content == "" {
		return errors.New("content of email is nil")
	}
	for _, v := range to {
		if !e.isEmail(v) {
			return errors.New("to email is error :" + v)
		}
	}
	return nil
}

func (e *Email) isEmail(to string) bool {
	if !strings.Contains(to, "@") {
		return false
	}
	emailRegexp, _ := regexp.Compile(emailRegexpStr)
	return emailRegexp.MatchString(to)
}

func (e *Email) addr() string {
	return fmt.Sprintf("%s:%s", e.config.Host, e.config.Port)
}

func (e *Email) plainAuth() smtp.Auth {
	return smtp.PlainAuth("", e.config.User, e.config.Password, e.config.Host)
}

func (e *Email) cloesClient(client *smtp.Client) {
	_ = client.Close()
}
