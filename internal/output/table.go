package output

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
)

type TableFormatter struct{}

func (f *TableFormatter) Format(w io.Writer, data interface{}, columns []Column) error {
	rows, err := toRows(data, columns)
	if err != nil {
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	// Header
	headers := make([]string, len(columns))
	for i, col := range columns {
		headers[i] = col.Header
	}
	fmt.Fprintln(tw, strings.Join(headers, "\t"))

	// Rows
	for _, row := range rows {
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}

	return tw.Flush()
}

func toRows(data interface{}, columns []Column) ([][]string, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var items []map[string]interface{}
	if err := json.Unmarshal(raw, &items); err != nil {
		var item map[string]interface{}
		if err := json.Unmarshal(raw, &item); err != nil {
			return nil, fmt.Errorf("data must be a JSON object or array")
		}
		items = []map[string]interface{}{item}
	}

	rows := make([][]string, 0, len(items))
	for _, item := range items {
		row := make([]string, len(columns))
		for i, col := range columns {
			val, ok := item[col.Field]
			if !ok {
				row[i] = ""
				continue
			}
			row[i] = formatValue(val)
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func formatValue(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Slice, reflect.Map:
			b, _ := json.Marshal(v)
			return string(b)
		default:
			return fmt.Sprintf("%v", v)
		}
	}
}
