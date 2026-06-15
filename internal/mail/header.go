package mail

import "strings"

// Header stores message headers in insertion order.
type Header struct {
	fields []headerField
}

type headerField struct {
	Name   string
	Values []string
}

// Add appends a header value.
func (h *Header) Add(name string, values ...string) error {
	if err := validateHeader(name, values...); err != nil {
		return err
	}
	h.fields = append(h.fields, headerField{Name: canonicalHeaderName(name), Values: append([]string(nil), values...)})
	return nil
}

// Set replaces all values for name while keeping deterministic order.
func (h *Header) Set(name string, values ...string) error {
	if err := validateHeader(name, values...); err != nil {
		return err
	}
	name = canonicalHeaderName(name)
	for i := range h.fields {
		if strings.EqualFold(h.fields[i].Name, name) {
			h.fields[i].Values = append([]string(nil), values...)
			return nil
		}
	}
	h.fields = append(h.fields, headerField{Name: name, Values: append([]string(nil), values...)})
	return nil
}

// Values returns all values for name.
func (h *Header) Values(name string) []string {
	for _, field := range h.fields {
		if strings.EqualFold(field.Name, name) {
			return append([]string(nil), field.Values...)
		}
	}
	return nil
}

func validateHeader(name string, values ...string) error {
	name = strings.TrimSpace(name)
	if name == "" || hasCRLF(name) || strings.Contains(name, ":") {
		return ErrInvalidHeader
	}
	for _, value := range values {
		if hasCRLF(value) {
			return ErrInvalidHeader
		}
	}
	return nil
}

func canonicalHeaderName(name string) string { return strings.TrimSpace(name) }

func hasCRLF(value string) bool { return strings.ContainsAny(value, "\r\n") }

func writeHeaderLine(b *strings.Builder, name string, values ...string) {
	b.WriteString(name)
	b.WriteString(":")
	if len(values) > 0 {
		b.WriteString(" ")
		b.WriteString(strings.Join(values, ", "))
	}
	b.WriteString("\r\n")
}
