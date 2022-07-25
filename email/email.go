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

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:38
 * @description: Email Client
 */
type Email struct {
	host     string
	port     string
	user     string
	passWord string
	from     string
}

var (
	emailRegexpStr = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
)

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:39
 * @params: host, port, user, password, from string
 * @return: *Email, error
 * @description: 新建Email Client
 */
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
	return &Email{host: host, port: port, user: user, passWord: password, from: from}, nil
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:39
 * @params: to []string, subject, content string
 * @return: error
 * @description: 发送普通邮件
 */
func (e *Email) Send(to []string, subject, content string) error {
	err := e.check(to, subject, content)
	if err != nil {
		return err
	}
	return e.sendMail(to, subject, content)
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:39
 * @params: to []string, subject, content string
 * @return: error
 * @description: 发送TLS加密邮件
 */
func (e *Email) SendTLS(to []string, subject, content string) error {
	err := e.check(to, subject, content)
	if err != nil {
		return err
	}
	return e.sendTLSMail(to, subject, content)
}

func (e *Email) sendMail(to []string, subject, content string) error {
	email := email.NewEmail()
	email.From = e.from
	email.To = to
	email.Subject = subject
	email.Text = []byte(content)
	err := email.Send(fmt.Sprintf("%s:%s", e.host, e.port),
		smtp.PlainAuth("", e.user, e.passWord, e.host))
	if err != nil {
		return err
	}
	return nil
}

func (e *Email) sendTLSMail(to []string, subject, content string) error {
	header := make(map[string]any)
	header["From"] = e.from
	header["Subject"] = subject
	header["Content-Type"] = "text/html; charset=UTF-8"
	body := content
	sendMsg := ""
	for k, v := range header {
		sendMsg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	sendMsg += "\r\n" + body
	err := e.sendMailUsingTLS(
		fmt.Sprintf("%s:%s", e.host, e.port),
		smtp.PlainAuth("", e.user, e.passWord, e.host),
		e.user,
		to,
		[]byte(sendMsg),
	)
	if err != nil {
		return err
	}
	return nil
}

func (e *Email) sendMailUsingTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) (err error) {
	client, err := e.dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()
	if auth != nil {
		if ok, _ := client.Extension("AUTH"); ok {
			if err = client.Auth(auth); err != nil {
				return err
			}
		}
	}
	if err = client.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := client.Data()
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
