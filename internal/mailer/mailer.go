package mailer

import "embed"

const (
	fromName            = "Go social media"
	maxTries            = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

//go:embed "templates"
var fs embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, issandbox bool) (int, error)
}
