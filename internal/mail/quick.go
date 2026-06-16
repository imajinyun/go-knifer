package mail

import "context"

// QuickOption customizes account-based quick send helpers.
type QuickOption func(*quickConfig) error

type quickConfig struct {
	messageOptions []MessageOption
	clientOptions  []ClientOption
}

// WithQuickMessageOptions appends message options used by QuickSend and account helpers.
func WithQuickMessageOptions(opts ...MessageOption) QuickOption {
	return func(q *quickConfig) error {
		q.messageOptions = append(q.messageOptions, opts...)
		return nil
	}
}

// WithQuickClientOptions appends client options used by QuickSend and account helpers.
func WithQuickClientOptions(opts ...ClientOption) QuickOption {
	return func(q *quickConfig) error {
		q.clientOptions = append(q.clientOptions, opts...)
		return nil
	}
}

// QuickSend creates and sends a message using account defaults plus quick options.
func QuickSend(ctx context.Context, account Account, opts ...QuickOption) error {
	config, err := newQuickConfig(opts...)
	if err != nil {
		return err
	}
	messageOpts, err := account.messageOptions(config.messageOptions...)
	if err != nil {
		return err
	}
	message, err := NewMessage(messageOpts...)
	if err != nil {
		return err
	}
	return sendAccountMessage(ctx, account, message, config.clientOptions...)
}

// SendAccountText creates and sends a plain text message using account defaults.
func SendAccountText(
	ctx context.Context,
	account Account,
	to []string,
	subject string,
	text string,
	opts ...QuickOption,
) error {
	messageOptions := []MessageOption{
		WithTo(to...),
		WithSubject(subject),
		WithText(text),
	}
	return sendAccountBody(ctx, account, messageOptions, opts...)
}

// SendAccountHTML creates and sends an HTML message using account defaults.
func SendAccountHTML(
	ctx context.Context,
	account Account,
	to []string,
	subject string,
	html string,
	opts ...QuickOption,
) error {
	messageOptions := []MessageOption{
		WithTo(to...),
		WithSubject(subject),
		WithHTML(html),
	}
	return sendAccountBody(ctx, account, messageOptions, opts...)
}

func sendAccountBody(
	ctx context.Context,
	account Account,
	messageOptions []MessageOption,
	opts ...QuickOption,
) error {
	quickOpts := make([]QuickOption, 0, 1+len(opts))
	quickOpts = append(quickOpts, WithQuickMessageOptions(messageOptions...))
	quickOpts = append(quickOpts, opts...)
	return QuickSend(ctx, account, quickOpts...)
}

func sendAccountMessage(ctx context.Context, account Account, message *Message, opts ...ClientOption) error {
	client, err := NewClient(account.Host, account.Port, account.clientOptions(opts...)...)
	if err != nil {
		return err
	}
	return client.Send(ctx, message)
}

func newQuickConfig(opts ...QuickOption) (quickConfig, error) {
	config := quickConfig{
		messageOptions: []MessageOption{},
		clientOptions:  []ClientOption{},
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(&config); err != nil {
			return quickConfig{}, err
		}
	}
	return config, nil
}
