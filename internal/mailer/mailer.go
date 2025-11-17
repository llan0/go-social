package mailer

import "embed"

const (
	FromName            = "GoSocial"
	MaxRetries          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, mail string, data any, isSandbox bool) (string, error)
}
