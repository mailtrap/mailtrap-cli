package output

import (
	"encoding/json"
	"fmt"
	"io"
)

// Page describes where a paginated list response keeps its items, total
// count and next-page cursor.
type Page struct {
	Items      string   // key of the items array
	Total      string   // key of the total count; empty when the API returns none
	Cursor     []string // path to the next-page cursor
	CursorFlag string   // flag that takes the cursor on the next request
	NextArgs   string   // extra flags the next request repeats
}

// PrintPage prints a paginated list response. JSON output is the body as the
// API returned it, so scripts keep every field, the total and the cursor.
// Table and text render the items followed by the total and next-page flag.
func PrintPage(w io.Writer, format Format, body json.RawMessage, page Page, columns []Column) error {
	if format == FormatJSON {
		return Print(w, format, body, columns)
	}

	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(body, &envelope); err != nil {
		return fmt.Errorf("data must be a JSON object")
	}
	if err := Print(w, format, envelope[page.Items], columns); err != nil {
		return err
	}

	var footer []string
	if total := lookup(body, page.Total); total != "" {
		footer = append(footer, "Total: "+total)
	}
	if cursor := lookup(body, page.Cursor...); cursor != "" {
		next := fmt.Sprintf("Next page: --%s %s", page.CursorFlag, cursor)
		if page.NextArgs != "" {
			next += " " + page.NextArgs
		}
		footer = append(footer, next)
	}
	if len(footer) > 0 {
		fmt.Fprintln(w)
		for _, line := range footer {
			fmt.Fprintln(w, line)
		}
	}
	return nil
}

// lookup returns the scalar at path in body, or "" when the path is empty,
// unset or missing from the response.
func lookup(body json.RawMessage, path ...string) string {
	if len(path) == 0 || path[0] == "" {
		return ""
	}
	var v interface{}
	if json.Unmarshal(body, &v) != nil {
		return ""
	}
	for _, key := range path {
		m, ok := v.(map[string]interface{})
		if !ok {
			return ""
		}
		v = m[key]
	}
	return formatValue(v)
}
