package vxml

import (
	"io"
	"io/fs"

	xmlimpl "github.com/imajinyun/go-knifer/internal/xml"
)

// WithNamespaceAware controls whether parsed element names keep namespace URIs.
func WithNamespaceAware(b bool) ParseOption { return xmlimpl.WithNamespaceAware(b) }

// WithStrict controls XML decoder strict mode.
func WithStrict(b bool) ParseOption { return xmlimpl.WithStrict(b) }

// WithCharsetReader sets the charset reader used by the XML decoder.
func WithCharsetReader(reader func(charset string, input io.Reader) (io.Reader, error)) ParseOption {
	return xmlimpl.WithCharsetReader(reader)
}

// WithEntity sets custom XML decoder entity replacements.
func WithEntity(entity map[string]string) ParseOption { return xmlimpl.WithEntity(entity) }

// WithMaxBytes bounds XML input read from readers and files. Non-positive values mean unlimited.
func WithMaxBytes(maxBytes int64) ParseOption { return xmlimpl.WithMaxBytes(maxBytes) }

// WithCharset sets the XML declaration charset.
func WithCharset(s string) WriteOption { return xmlimpl.WithCharset(s) }

// WithIndent sets the indentation width in spaces (0 disables pretty printing).
func WithIndent(n int) WriteOption { return xmlimpl.WithIndent(n) }

// WithPretty enables pretty printing with the default indentation.
func WithPretty() WriteOption { return xmlimpl.WithPretty() }

// WithOmitDeclaration controls whether the <?xml ... ?> prolog is emitted.
func WithOmitDeclaration(b bool) WriteOption { return xmlimpl.WithOmitDeclaration(b) }

// WithIgnoreNullFields skips struct fields whose value is a typed nil.
func WithIgnoreNullFields(b bool) WriteOption { return xmlimpl.WithIgnoreNullFields(b) }

// WithRootName overrides the synthesized root element name for MarshalMap / MarshalBean.
func WithRootName(s string) WriteOption { return xmlimpl.WithRootName(s) }

// WithNamespace sets the xmlns attribute on the synthesized root element.
func WithNamespace(s string) WriteOption { return xmlimpl.WithNamespace(s) }

// WithFilePerm sets the file permission used by WriteFile.
func WithFilePerm(perm fs.FileMode) WriteOption { return xmlimpl.WithFilePerm(perm) }

// WithDirPerm sets the parent-directory permission used by WriteFile.
func WithDirPerm(perm fs.FileMode) WriteOption { return xmlimpl.WithDirPerm(perm) }

// WithOverwrite controls whether WriteFile may replace an existing file.
func WithOverwrite(overwrite bool) WriteOption { return xmlimpl.WithOverwrite(overwrite) }

// WithCreateParents controls whether WriteFile creates parent directories.
func WithCreateParents(create bool) WriteOption { return xmlimpl.WithCreateParents(create) }
