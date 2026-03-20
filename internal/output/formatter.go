package output

import (
	"fmt"
	"io"
)

type Format string

const (
	FormatJSON  Format = "json"
	FormatTable Format = "table"
	FormatText  Format = "text"
)

type Column struct {
	Header string
	Field  string
}

type Formatter interface {
	Format(w io.Writer, data interface{}, columns []Column) error
}

func New(format Format) Formatter {
	switch format {
	case FormatJSON:
		return &JSONFormatter{}
	case FormatText:
		return &TextFormatter{}
	default:
		return &TableFormatter{}
	}
}

func Print(w io.Writer, format Format, data interface{}, columns []Column) error {
	return New(format).Format(w, data, columns)
}

func PrintRaw(w io.Writer, data []byte) error {
	_, err := fmt.Fprint(w, string(data))
	return err
}
