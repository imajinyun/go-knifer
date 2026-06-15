package mail

import (
	"errors"
	stdmail "net/mail"
	"strings"
)

// Address is an RFC 5322 mailbox address.
type Address struct {
	Name  string
	Email string
}

// ParseAddress parses a single mailbox address.
func ParseAddress(value string) (*Address, error) {
	if hasCRLF(value) {
		return nil, ErrInvalidAddress
	}
	addr, err := stdmail.ParseAddress(strings.TrimSpace(value))
	if err != nil {
		return nil, errors.Join(ErrInvalidAddress, err)
	}
	if err := validateEmail(addr.Address); err != nil {
		return nil, err
	}
	return &Address{Name: addr.Name, Email: addr.Address}, nil
}

// ParseAddressList parses a comma-separated mailbox address list.
func ParseAddressList(value string) ([]*Address, error) {
	if hasCRLF(value) {
		return nil, ErrInvalidAddress
	}
	parsed, err := stdmail.ParseAddressList(strings.TrimSpace(value))
	if err != nil {
		return nil, errors.Join(ErrInvalidAddress, err)
	}
	addresses := make([]*Address, 0, len(parsed))
	for _, addr := range parsed {
		if err := validateEmail(addr.Address); err != nil {
			return nil, err
		}
		addresses = append(addresses, &Address{Name: addr.Name, Email: addr.Address})
	}
	return addresses, nil
}

// NewAddress validates and returns a mailbox address.
func NewAddress(name, email string) (*Address, error) {
	if hasCRLF(name) || hasCRLF(email) {
		return nil, ErrInvalidAddress
	}
	if err := validateEmail(email); err != nil {
		return nil, err
	}
	return &Address{Name: name, Email: email}, nil
}

// String formats a for use in message headers.
func (a *Address) String() string {
	if a == nil {
		return ""
	}
	return (&stdmail.Address{Name: a.Name, Address: a.Email}).String()
}

func formatAddressList(addresses []*Address) string {
	parts := make([]string, 0, len(addresses))
	for _, addr := range addresses {
		if addr != nil {
			parts = append(parts, addr.String())
		}
	}
	return strings.Join(parts, ", ")
}

func validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" || hasCRLF(email) {
		return ErrInvalidAddress
	}
	addr, err := stdmail.ParseAddress(email)
	if err != nil || addr.Address != email {
		if err != nil {
			return errors.Join(ErrInvalidAddress, err)
		}
		return ErrInvalidAddress
	}
	return nil
}

func cloneAddresses(in []*Address) []*Address {
	if len(in) == 0 {
		return nil
	}
	out := make([]*Address, 0, len(in))
	for _, addr := range in {
		if addr == nil {
			continue
		}
		copyAddr := *addr
		out = append(out, &copyAddr)
	}
	return out
}
