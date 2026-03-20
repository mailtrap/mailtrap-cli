package output

import (
	"encoding/json"
	"fmt"
	"io"
)

type TextFormatter struct{}

func (f *TextFormatter) Format(w io.Writer, data interface{}, columns []Column) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Try as array
	var items []map[string]interface{}
	if err := json.Unmarshal(raw, &items); err != nil {
		// Single object
		var item map[string]interface{}
		if err := json.Unmarshal(raw, &item); err != nil {
			return fmt.Errorf("data must be a JSON object or array")
		}
		items = []map[string]interface{}{item}
	}

	for i, item := range items {
		if i > 0 {
			fmt.Fprintln(w, "---")
		}
		for _, col := range columns {
			val := item[col.Field]
			fmt.Fprintf(w, "%s: %s\n", col.Header, formatValue(val))
		}
	}
	return nil
}
