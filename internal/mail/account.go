package mail

import (
	"crypto/tls"
	"net/smtp"
	"strings"
	"time"
)

// Account stores SMTP server, authentication, and default sender settings.
type Account struct {
	Host           string
	Port           int
	Username       string
	Password       string
	Auth           smtp.Auth
	From           string
	FromName       string
	TLSConfig      *tls.Config
	TLSPolicy      TLSPolicy
	AllowPlainAuth bool
	Timeout        time.Duration
	LocalName      string
}

func (a Account) clientOptions(extra ...ClientOption) []ClientOption {
	opts := make([]ClientOption, 0, 8+len(extra))
	if a.Username != "" || a.Password != "" {
		opts = append(opts, WithAuth(a.Username, a.Password))
	}
	if a.Auth != nil {
		opts = append(opts, WithSMTPAuth(a.Auth))
	}
	if a.TLSConfig != nil {
		opts = append(opts, WithTLSConfig(a.TLSConfig))
	}
	if a.TLSPolicy != TLSPolicyUnknown {
		opts = append(opts, WithTLSPolicy(a.TLSPolicy))
	}
	if a.AllowPlainAuth {
		opts = append(opts, WithAllowPlainAuth(true))
	}
	if a.Timeout > 0 {
		opts = append(opts, WithTimeout(a.Timeout))
	}
	if a.LocalName != "" {
		opts = append(opts, WithLocalName(a.LocalName))
	}
	opts = append(opts, extra...)
	return opts
}

func (a Account) messageOptions(extra ...MessageOption) ([]MessageOption, error) {
	from := strings.TrimSpace(a.From)
	if from == "" {
		from = strings.TrimSpace(a.Username)
	}
	if from == "" {
		return nil, ErrMissingFrom
	}

	opts := make([]MessageOption, 0, 1+len(extra))
	if a.FromName == "" {
		opts = append(opts, WithFrom(from))
	} else {
		addr, err := NewAddress(a.FromName, from)
		if err != nil {
			return nil, err
		}
		opts = append(opts, WithFromAddress(addr))
	}
	opts = append(opts, extra...)
	return opts, nil
}
