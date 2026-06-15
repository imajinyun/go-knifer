package mail

// Charset identifies a MIME character set.
type Charset string

// Encoding identifies a Content-Transfer-Encoding value.
type Encoding string

// ContentType identifies a MIME media type.
type ContentType string

const (
	// CharsetUTF8 is the default charset used for text parts and encoded words.
	CharsetUTF8 Charset = "UTF-8"
	// CharsetASCII identifies US-ASCII content.
	CharsetASCII Charset = "US-ASCII"
)

const (
	// EncodingBase64 encodes content using RFC 2045 base64.
	EncodingBase64 Encoding = "base64"
	// EncodingQuotedPrintable encodes content using quoted-printable.
	EncodingQuotedPrintable Encoding = "quoted-printable"
	// Encoding7Bit marks content as 7bit.
	Encoding7Bit Encoding = "7bit"
	// Encoding8Bit marks content as 8bit.
	Encoding8Bit Encoding = "8bit"
)

const (
	// TypeTextPlain is plain text content.
	TypeTextPlain ContentType = "text/plain"
	// TypeTextHTML is HTML content.
	TypeTextHTML ContentType = "text/html"
	// TypeApplicationOctetStream is generic binary content.
	TypeApplicationOctetStream ContentType = "application/octet-stream"
)

func (c Charset) String() string { return string(c) }

func (e Encoding) String() string { return string(e) }

func (c ContentType) String() string { return string(c) }
