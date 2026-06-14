package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/wneessen/go-mail"
)

//go:embed templates
var emailTemp embed.FS // mail templates

type Mailer struct {
	Dialer *mail.Client
	Sender string
}

func New(host string, port int, username string, password string, sender string) (*Mailer, error) {
	client, err := mail.NewClient(host,
		mail.WithPort(port),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
	)
	if err != nil {
		return nil, err
	}

	return &Mailer{
		Dialer: client,
		Sender: sender,
	}, nil
}

func (m *Mailer) Send(recipient string, filename string, data any) error {
	t, err := template.New("email").ParseFS(emailTemp, "templates/"+filename)
	if err != nil {
		return err
	}

	var subject bytes.Buffer
	err = t.ExecuteTemplate(&subject, "subject", data)
	if err != nil {
		return err
	}

	var htmlBody bytes.Buffer
	err = t.ExecuteTemplate(&htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	var textBody bytes.Buffer
	err = t.ExecuteTemplate(&textBody, "plainBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMsg()
	err = msg.To(recipient)
	if err != nil {
		return err
	}

	err = msg.From(m.Sender)
	if err != nil {
		return err
	}

	msg.Subject(subject.String())
	msg.SetBodyString(mail.TypeTextHTML, htmlBody.String())
	msg.AddAlternativeString(mail.TypeTextPlain, textBody.String())

	for range 3 {
		err = m.Dialer.DialAndSend(msg)
		if err == nil {
			return nil
		}

		time.Sleep(time.Millisecond * 200)
	}

	return err
}
