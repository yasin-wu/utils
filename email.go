package utils

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/jordan-wright/email"
	"net"
	"net/smtp"
	"regexp"
	"strings"
)

type Email struct {
	From       string
	To         string
	Subject    string
	Content    string
	SMTPServer *SMTPServer
}

type SMTPServer struct {
	Host     string
	Port     string
	User     string
	Password string
}

var (
	emailRegexpStr = `^[_a-z0-9-]+(\.[_a-z0-9-]+)*@[a-z0-9-]+(\.[a-z0-9-]+)*(\.[a-z]{2,})$`
)

func (this *Email) Send() error {
	err := sendMail(
		this.SMTPServer.Host,
		this.SMTPServer.Port,
		this.SMTPServer.User,
		this.SMTPServer.Password,
		this.To,
		this.Subject,
		this.Content)
	if err != nil {
		return errors.New("send mail failed! err:" + err.Error())
	}
	return nil
}

func (this *Email) SendTLS() error {
	err := sendTLSMail(
		this.SMTPServer.Host,
		this.SMTPServer.Port,
		this.SMTPServer.User,
		this.SMTPServer.Password,
		this.To,
		this.Subject,
		this.Content)
	if err != nil {
		return errors.New("send tlsmail failed! err:" + err.Error())
	}
	return nil
}

func (this *Email) IsEmail() bool {
	if !strings.Contains(this.To, "@") {
		return false
	}
	emailRegexp, _ := regexp.Compile(emailRegexpStr)
	return emailRegexp.MatchString(this.To)
}

func sendMail(host, port, user, password, to, subject, content string) error {
	e := email.NewEmail()
	e.From = user
	e.To = []string{to}
	e.Subject = subject
	e.Text = []byte(content)
	err := e.Send(host+":"+port, smtp.PlainAuth("", user, password, host))
	if err != nil {
		return err
	}
	return nil
}

func sendTLSMail(host, port, user, password, to, subject, content string) error {
	header := make(map[string]string)
	header["From"] = user
	header["To"] = to
	header["Subject"] = subject
	header["Content-Type"] = "text/html; charset=UTF-8"
	body := content
	sendMsg := ""
	for k, v := range header {
		sendMsg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	sendMsg += "\r\n" + body
	auth := smtp.PlainAuth(
		"",
		user,
		password,
		host,
	)
	err := sendMailUsingTLS(
		fmt.Sprintf("%s:%s", host, port),
		auth,
		user,
		[]string{to},
		[]byte(sendMsg),
	)
	if err != nil {
		return err
	}
	return nil
}

func sendMailUsingTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) (err error) {
	c, err := dial(addr)
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

func dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}
