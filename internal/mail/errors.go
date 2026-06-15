package mail

import "errors"

var (
	// ErrInvalidAddress is returned when an email address cannot be parsed or validated.
	ErrInvalidAddress = errors.New("mail: invalid address")
	// ErrInvalidHeader is returned when a header name or value is not safe for SMTP/MIME output.
	ErrInvalidHeader = errors.New("mail: invalid header")
	// ErrMissingFrom is returned when a message has no From address.
	ErrMissingFrom = errors.New("mail: missing from address")
	// ErrMissingRecipient is returned when a message has no To, Cc, or Bcc recipient.
	ErrMissingRecipient = errors.New("mail: missing recipient")
	// ErrMissingBody is returned when a message has no body content.
	ErrMissingBody = errors.New("mail: missing body")
	// ErrTLSRequired is returned when the configured security policy requires TLS but TLS is unavailable.
	ErrTLSRequired = errors.New("mail: tls required")
	// ErrPlainAuth is returned when SMTP AUTH would be sent over a plaintext connection.
	ErrPlainAuth = errors.New("mail: plaintext auth disabled")
	// ErrAttachmentTooLarge is returned when an attachment exceeds the configured size limit.
	ErrAttachmentTooLarge = errors.New("mail: attachment too large")
)
