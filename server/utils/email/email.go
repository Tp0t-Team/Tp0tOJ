package email

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/jordan-wright/email"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"server/utils/configure"
	"strings"
)

func SendMail(address string, subject string, content []byte) error {

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	mail := email.NewEmail()
	err := checkmail.ValidateFormat(configure.Configure.Email.Username)
	if err != nil {
		log.Println("mail check failed ", err)
	}
	mail.From = strings.Split(configure.Configure.Email.Username, "@")[0] + fmt.Sprintf("<%s>", configure.Configure.Email.Username)
	mail.To = []string{address}
	mail.Subject = subject
	mail.Text = content

	return mail.Send(fmt.Sprintf("%s:25", configure.Configure.Email.Host), smtp.PlainAuth("", configure.Configure.Email.Username, configure.Configure.Email.Password, configure.Configure.Email.Host))
}

//go:embed reset.html
var embedResetTemplate string

var resetTemplate *template.Template

const resetTemplatePath = "resources/emails/reset.html"

func loadTemplate(file string, defaultText string) (*template.Template, error) {
	if _, err := os.Stat(file); err == nil {
		text, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		return template.New("reset").Parse(string(text))
	} else {
		return template.New("reset").Parse(defaultText)
	}
}

func init() {
	var err error
	resetTemplate, err = loadTemplate(resetTemplatePath, embedResetTemplate)
	if err != nil {
		log.Panicln(err)
	}
}

type ResetInfo struct {
	Username string
	Url      string
}

func RenderResetEmail(info ResetInfo) ([]byte, error) {
	var ret bytes.Buffer
	if err := resetTemplate.Execute(&ret, info); err != nil {
		return nil, err
	}
	return ret.Bytes(), nil
}

func RenderWelcomeEmail(info ResetInfo) ([]byte, error) {
	// TODO:
	return RenderResetEmail(info)
}
