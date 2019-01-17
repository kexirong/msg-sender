package email

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"
)

type SMTP struct {
	address,
	username,
	password,
	authtype string
}

func New(address, username, password, authtype string) *SMTP {
	return &SMTP{
		address:  address,
		username: username,
		password: password,
		authtype: authtype,
	}
}

func (self *SMTP) SendMail(from string, to []string, subject, body string, contentType ...string) error {

	tos := strings.Join(to, ";")
	addrArr := strings.Split(self.address, ":")
	if len(addrArr) != 2 {
		return fmt.Errorf("address format error")
	}

	b64 := base64.StdEncoding // base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

	header := make(map[string]string)

	header["From"] = self.username
	if strings.HasSuffix(from, "staff.qkagame.com") {
		header["From"] = from
	}
	header["To"] = tos
	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", b64.EncodeToString([]byte(subject)))
	header["MIME-Version"] = "1.0"

	ctType := "text/plain; charset=UTF-8"
	if len(contentType) > 0 && contentType[0] == "html" {
		ctType = "text/html; charset=UTF-8"
	}
	header["Content-Type"] = ctType
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + b64.EncodeToString([]byte(body))
	var auth smtp.Auth
	switch self.authtype {
	case "LOGIN":
		auth = LoginAuth(self.username, self.password)
	case "CRAM-MD5":
		auth = smtp.CRAMMD5Auth(self.username, self.password)
	default:
		auth = smtp.PlainAuth("", self.username, self.password, addrArr[0])

	}

	return smtp.SendMail(self.address, auth, self.username, to, []byte(message))
}
