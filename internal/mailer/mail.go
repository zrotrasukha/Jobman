package mailer

import (
	"bytes"
	"embed"
	"html/template"

	"github.com/wneessen/go-mail"
)

//go:embed templates
var emailTemp embed.FS

type Mailer struct {
	Dialer *mail.Client
	Sender string
}

func New(host string, port int, username string, password string, sender string) (*Mailer, error) {
	client, err := mail.NewClient(host,
		mail.WithPort(port),
		mail.WithUsername(username),
		mail.WithPassword(password),
	)
	if err != nil {
		return nil, err
	}

	return &Mailer{
		Dialer: client,
		Sender: sender,
	}, nil
}

func (m *Mailer) Send(recipient string, sender string, data any) error {
	t, err := template.New("email").ParseFS(emailTemp, "templates/email.html")
	if err != nil {
		return err
	}

	var subject bytes.Buffer
	err = t.ExecuteTemplate(&subject, "subject", nil)
	if err != nil {
		return err
	}

	var htmlBody bytes.Buffer
	err = t.ExecuteTemplate(&htmlBody, "html", nil)
	if err != nil {
		return err
	}

	var textBody bytes.Buffer
	err = t.ExecuteTemplate(&textBody, "text", nil)
	if err != nil {
		return err
	}

	msg := mail.NewMsg()
	err = msg.To(recipient)
	if err != nil {
		return err
	}

	err = msg.From(sender)
	if err != nil {
		return err
	}

	msg.Subject(subject.String())
	msg.SetBodyString(mail.TypeTextHTML, htmlBody.String())
	msg.AddAlternativeString(mail.TypeTextPlain, textBody.String())

	err = m.Dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
