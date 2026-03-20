package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// 1. New() returns correct formatter type for each format
// ---------------------------------------------------------------------------

func TestNew_JSONFormat(t *testing.T) {
	f := New(FormatJSON)
	if _, ok := f.(*JSONFormatter); !ok {
		t.Errorf("New(FormatJSON) returned %T, want *JSONFormatter", f)
	}
}

func TestNew_TableFormat(t *testing.T) {
	f := New(FormatTable)
	if _, ok := f.(*TableFormatter); !ok {
		t.Errorf("New(FormatTable) returned %T, want *TableFormatter", f)
	}
}

func TestNew_TextFormat(t *testing.T) {
	f := New(FormatText)
	if _, ok := f.(*TextFormatter); !ok {
		t.Errorf("New(FormatText) returned %T, want *TextFormatter", f)
	}
}

func TestNew_UnknownFormatDefaultsToTable(t *testing.T) {
	f := New(Format("yaml"))
	if _, ok := f.(*TableFormatter); !ok {
		t.Errorf("New(unknown) returned %T, want *TableFormatter (default)", f)
	}
}

// ---------------------------------------------------------------------------
// 2. JSONFormatter – single object
// ---------------------------------------------------------------------------

func TestJSONFormatter_SingleObject(t *testing.T) {
	var buf bytes.Buffer
	f := &JSONFormatter{}
	data := map[string]interface{}{
		"id":   1,
		"name": "Alice",
	}
	cols := []Column{{Header: "ID", Field: "id"}, {Header: "Name", Field: "name"}}

	if err := f.Format(&buf, data, cols); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if got["name"] != "Alice" {
		t.Errorf("name = %v, want Alice", got["name"])
	}
}

// ---------------------------------------------------------------------------
// 3. JSONFormatter – array of objects
// ---------------------------------------------------------------------------

func TestJSONFormatter_Array(t *testing.T) {
	var buf bytes.Buffer
	f := &JSONFormatter{}
	data := []map[string]interface{}{
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
	}
	cols := []Column{{Header: "ID", Field: "id"}}

	if err := f.Format(&buf, data, cols); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON array: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("got %d items, want 2", len(got))
	}
}

// ---------------------------------------------------------------------------
// 4. TableFormatter – single object with columns
// ---------------------------------------------------------------------------

func TestTableFormatter_SingleObject(t *testing.T) {
	var buf bytes.Buffer
	f := &TableFormatter{}
	data := map[string]interface{}{
		"id":   42,
		"name": "Widget",
	}
	cols := []Column{
		{Header: "ID", Field: "id"},
		{Header: "NAME", Field: "name"},
	}

	if err := f.Format(&buf, data, cols); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected at least 2 lines (header + row), got %d:\n%s", len(lines), out)
	}
	if !strings.Contains(lines[0], "ID") || !strings.Contains(lines[0], "NAME") {
		t.Errorf("header line missing expected columns: %q", lines[0])
	}
	if !strings.Contains(lines[1], "42") || !strings.Contains(lines[1], "Widget") {
		t.Errorf("data row missing expected values: %q", lines[1])
	}
}

// ---------------------------------------------------------------------------
// 5. TableFormatter – array with correct headers and rows
// ---------------------------------------------------------------------------

func TestTableFormatter_Array(t *testing.T) {
	var buf bytes.Buffer
	f := &TableFormatter{}
	data := []map[string]interface{}{
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
	}
	cols := []Column{
		{Header: "ID", Field: "id"},
		{Header: "NAME", Field: "name"},
	}

	if err := f.Format(&buf, data, cols); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// 1 header + 2 data rows
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d:\n%s", len(lines), buf.String())
	}
	if !strings.Contains(lines[1], "Alice") {
		t.Errorf("first data row missing Alice: %q", lines[1])
	}
	if !strings.Contains(lines[2], "Bob") {
		t.Errorf("second data row missing Bob: %q", lines[2])
	}
}

// ---------------------------------------------------------------------------
// 6. TableFormatter – handles nil/empty/missing values
// ---------------------------------------------------------------------------

func TestTableFormatter_MissingAndNilValues(t *testing.T) {
	var buf bytes.Buffer
	f := &TableFormatter{}
	data := map[string]interface{}{
		"id":   1,
		"name": nil,
		// "email" is missing entirely
	}
	cols := []Column{
		{Header: "ID", Field: "id"},
		{Header: "NAME", Field: "name"},
		{Header: "EMAIL", Field: "email"},
	}

	if err := f.Format(&buf, data, cols); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	// nil and missing fields should produce empty strings, so no panic/error
	if !strings.Contains(out, "ID") {
		t.Errorf("output missing header: %s", out)
	}
}

// ---------------------------------------------------------------------------
// 7. TextFormatter – single object as key-value pairs
// ---------------------------------------------------------------------------

func TestTextFormatter_SingleObject(t *testing.T) {
	var buf bytes.Buffer
	f := &TextFormatter{}
	data := map[string]interface{}{
		"id":   10,
		"name": "Test",
	}
	cols := []Column{
		{Header: "ID", Field: "id"},
		{Header: "Name", Field: "name"},
	}

	if err := f.Format(&buf, data, cols); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ID: 10") {
		t.Errorf("expected 'ID: 10' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Name: Test") {
		t.Errorf("expected 'Name: Test' in output, got:\n%s", out)
	}
	if strings.Contains(out, "---") {
		t.Errorf("single object should not contain separator '---'")
	}
}

// ---------------------------------------------------------------------------
// 8. TextFormatter – array with separators
// ---------------------------------------------------------------------------

func TestTextFormatter_ArrayWithSeparators(t *testing.T) {
	var buf bytes.Buffer
	f := &TextFormatter{}
	data := []map[string]interface{}{
		{"id": 1, "name": "Alpha"},
		{"id": 2, "name": "Beta"},
	}
	cols := []Column{
		{Header: "ID", Field: "id"},
		{Header: "Name", Field: "name"},
	}

	if err := f.Format(&buf, data, cols); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "---") {
		t.Errorf("expected separator '---' between items, got:\n%s", out)
	}
	if strings.Count(out, "---") != 1 {
		t.Errorf("expected exactly 1 separator for 2 items, got %d", strings.Count(out, "---"))
	}
	if !strings.Contains(out, "Name: Alpha") {
		t.Errorf("missing first item data in output")
	}
	if !strings.Contains(out, "Name: Beta") {
		t.Errorf("missing second item data in output")
	}
}

// ---------------------------------------------------------------------------
// 9. Print() convenience function
// ---------------------------------------------------------------------------

func TestPrint_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]interface{}{"key": "value"}
	cols := []Column{{Header: "Key", Field: "key"}}

	if err := Print(&buf, FormatJSON, data, cols); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("Print with FormatJSON did not produce valid JSON: %v", err)
	}
	if got["key"] != "value" {
		t.Errorf("key = %v, want value", got["key"])
	}
}

func TestPrint_TableFormat(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]interface{}{"x": "hello"}
	cols := []Column{{Header: "X", Field: "x"}}

	if err := Print(&buf, FormatTable, data, cols); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "hello") {
		t.Errorf("table output missing value 'hello': %s", buf.String())
	}
}

func TestPrint_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]interface{}{"x": "world"}
	cols := []Column{{Header: "X", Field: "x"}}

	if err := Print(&buf, FormatText, data, cols); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "X: world") {
		t.Errorf("text output missing 'X: world': %s", buf.String())
	}
}

// ---------------------------------------------------------------------------
// 10. PrintRaw() outputs raw bytes
// ---------------------------------------------------------------------------

func TestPrintRaw(t *testing.T) {
	var buf bytes.Buffer
	raw := []byte("raw output content\nline two")

	if err := PrintRaw(&buf, raw); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != string(raw) {
		t.Errorf("PrintRaw output = %q, want %q", buf.String(), string(raw))
	}
}

func TestPrintRaw_EmptyBytes(t *testing.T) {
	var buf bytes.Buffer
	if err := PrintRaw(&buf, []byte{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != "" {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

// ---------------------------------------------------------------------------
// 11. formatValue – handles nil, string, int, float, slice, map
// ---------------------------------------------------------------------------

func TestFormatValue_Nil(t *testing.T) {
	if got := formatValue(nil); got != "" {
		t.Errorf("formatValue(nil) = %q, want empty string", got)
	}
}

func TestFormatValue_String(t *testing.T) {
	if got := formatValue("hello"); got != "hello" {
		t.Errorf("formatValue(string) = %q, want %q", got, "hello")
	}
}

func TestFormatValue_Int(t *testing.T) {
	got := formatValue(42)
	if got != "42" {
		t.Errorf("formatValue(42) = %q, want %q", got, "42")
	}
}

func TestFormatValue_Float(t *testing.T) {
	got := formatValue(3.14)
	if !strings.Contains(got, "3.14") {
		t.Errorf("formatValue(3.14) = %q, expected to contain '3.14'", got)
	}
}

func TestFormatValue_Bool(t *testing.T) {
	if got := formatValue(true); got != "true" {
		t.Errorf("formatValue(true) = %q, want %q", got, "true")
	}
}

func TestFormatValue_Slice(t *testing.T) {
	input := []interface{}{"a", "b"}
	got := formatValue(input)
	// Should be JSON-encoded
	if got != `["a","b"]` {
		t.Errorf("formatValue(slice) = %q, want %q", got, `["a","b"]`)
	}
}

func TestFormatValue_Map(t *testing.T) {
	input := map[string]interface{}{"k": "v"}
	got := formatValue(input)
	if got != `{"k":"v"}` {
		t.Errorf("formatValue(map) = %q, want %q", got, `{"k":"v"}`)
	}
}

// ---------------------------------------------------------------------------
// 12. toRows – handles invalid data gracefully
// ---------------------------------------------------------------------------

func TestToRows_InvalidData(t *testing.T) {
	cols := []Column{{Header: "X", Field: "x"}}
	// A plain string cannot be unmarshalled as object or array of objects
	_, err := toRows("not a valid object", cols)
	if err == nil {
		t.Error("expected error for invalid data, got nil")
	}
	if !strings.Contains(err.Error(), "data must be a JSON object or array") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestToRows_SingleObject(t *testing.T) {
	cols := []Column{
		{Header: "ID", Field: "id"},
		{Header: "Name", Field: "name"},
	}
	data := map[string]interface{}{"id": 1, "name": "Alice"}
	rows, err := toRows(data, cols)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if rows[0][1] != "Alice" {
		t.Errorf("row[0][1] = %q, want Alice", rows[0][1])
	}
}

func TestToRows_EmptyArray(t *testing.T) {
	cols := []Column{{Header: "X", Field: "x"}}
	data := []map[string]interface{}{}
	rows, err := toRows(data, cols)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rows) != 0 {
		t.Errorf("expected 0 rows for empty array, got %d", len(rows))
	}
}

func TestToRows_MissingField(t *testing.T) {
	cols := []Column{{Header: "X", Field: "nonexistent"}}
	data := map[string]interface{}{"id": 1}
	rows, err := toRows(data, cols)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rows[0][0] != "" {
		t.Errorf("expected empty string for missing field, got %q", rows[0][0])
	}
}

// ---------------------------------------------------------------------------
// Additional: struct input (common real-world usage)
// ---------------------------------------------------------------------------

type sampleItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestTableFormatter_StructInput(t *testing.T) {
	var buf bytes.Buffer
	f := &TableFormatter{}
	data := sampleItem{ID: 7, Name: "Gadget"}
	cols := []Column{
		{Header: "ID", Field: "id"},
		{Header: "NAME", Field: "name"},
	}

	if err := f.Format(&buf, data, cols); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "7") || !strings.Contains(out, "Gadget") {
		t.Errorf("struct data not rendered correctly: %s", out)
	}
}

func TestTextFormatter_StructSlice(t *testing.T) {
	var buf bytes.Buffer
	f := &TextFormatter{}
	data := []sampleItem{
		{ID: 1, Name: "A"},
		{ID: 2, Name: "B"},
	}
	cols := []Column{
		{Header: "ID", Field: "id"},
		{Header: "Name", Field: "name"},
	}

	if err := f.Format(&buf, data, cols); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "---") {
		t.Errorf("expected separator between struct items")
	}
	if !strings.Contains(out, "Name: A") || !strings.Contains(out, "Name: B") {
		t.Errorf("struct slice data not rendered correctly: %s", out)
	}
}

func TestJSONFormatter_StructSlice(t *testing.T) {
	var buf bytes.Buffer
	f := &JSONFormatter{}
	data := []sampleItem{
		{ID: 1, Name: "X"},
	}
	cols := []Column{}

	if err := f.Format(&buf, data, cols); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	if got[0]["name"] != "X" {
		t.Errorf("expected name=X, got %v", got[0]["name"])
	}
}

// ---------------------------------------------------------------------------
// Edge case: TextFormatter with invalid data
// ---------------------------------------------------------------------------

func TestTextFormatter_InvalidData(t *testing.T) {
	var buf bytes.Buffer
	f := &TextFormatter{}
	cols := []Column{{Header: "X", Field: "x"}}

	err := f.Format(&buf, "plain string", cols)
	if err == nil {
		t.Error("expected error for invalid data, got nil")
	}
}

func TestTableFormatter_InvalidData(t *testing.T) {
	var buf bytes.Buffer
	f := &TableFormatter{}
	cols := []Column{{Header: "X", Field: "x"}}

	err := f.Format(&buf, "plain string", cols)
	if err == nil {
		t.Error("expected error for invalid data, got nil")
	}
}
