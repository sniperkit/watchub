package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/matcornic/hermes"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sniperkit/watchub/config"
	"github.com/sniperkit/watchub/shared/dto"
)

var emailConfig = hermes.Hermes{
	Product: hermes.Product{
		Name:      "Watchub",
		Link:      "http://watchub.pw",
		Logo:      "https://raw.githubusercontent.com/caarlos0/watchub/master/static/apple-touch-icon-144x144.png",
		Copyright: "Copyright © 2016-2017 Watchub.",
	},
	Theme: new(hermes.Flat),
}

var welcomeIntro = []string{
	"Welcome to Watchub!",
	"We're very excited to have you on board.",
}

var changesIntro = []string{
	"Here is what changed in your account recently:",
}

func New(config config.Config) *MailSvc {
	return &MailSvc{
		hermes:  emailConfig,
		changes: template.Must(template.ParseFiles("static/mail/changes.md")),
		welcome: template.Must(template.ParseFiles("static/mail/welcome.md")),
		config:  config,
	}
}

type MailSvc struct {
	hermes  hermes.Hermes
	config  config.Config
	changes *template.Template
	welcome *template.Template
}

func (s *MailSvc) SendWelcome(data dto.WelcomeEmailData) {
	html, err := s.generate(data.Login, data, s.welcome, welcomeIntro)
	if err != nil {
		log.WithError(err).Error("failed to generate welcome email")
		return
	}
	s.send(data.Login, data.Email, "Welcome to Watchub!", html)
}

func (s *MailSvc) SendChanges(data dto.ChangesEmailData) {
	html, err := s.generate(data.Login, data, s.changes, changesIntro)
	if err != nil {
		log.WithError(err).Error("failed to generate changes email")
		return
	}
	s.send(data.Login, data.Email, "Your report from Watchub!", html)
}

func (s *MailSvc) generate(login string, data interface{}, tmpl *template.Template, intros []string) (string, error) {
	var wr bytes.Buffer
	if err := tmpl.Execute(&wr, data); err != nil {
		return "", err
	}
	return s.hermes.GenerateHTML(
		hermes.Email{
			Body: hermes.Body{
				Name:   login,
				Intros: intros,
				FreeMarkdown: hermes.Markdown(
					strings.Join(
						[]string{
							wr.String(),
							"\n\n",
							"We will continue to watch for changes and let you know!",
							"\n\n---\n\n",
							"<small>",
							`Liking our service? Maybe you'll consider [make a donation](http://watchub.pw/donate).`,
							fmt.Sprintf(
								`You might also want to change [your settings](%s).`,
								"https://github.com/settings/connections/applications/"+s.config.ClientID,
							),
							"</small>",
						},
						" ",
					),
				),
			},
		},
	)
}

func (s *MailSvc) send(name, email, subject, html string) {
	var log = log.WithField("email", email)
	var from = mail.NewEmail("Watchub", "noreply@watchub.pw")
	var to = mail.NewEmail(name, email)
	var request = sendgrid.GetRequest(
		s.config.SendgridAPIKey,
		"/v3/mail/send",
		"https://api.sendgrid.com",
	)
	// prevent grouping in gmail
	request.Headers["X-Entity-Ref-ID"] = "watchub-" + time.Now().String()
	request.Method = "POST"
	request.Body = mail.GetRequestBody(
		mail.NewV3MailInit(
			from,
			subject,
			to,
			mail.NewContent("text/html", html),
		),
	)
	resp, err := sendgrid.API(request)
	if err != nil {
		log.WithError(err).WithField("resp", resp).Error("failed to send email")
		return
	}
	log.WithField("status", resp.StatusCode).Info("email request posted")
}
