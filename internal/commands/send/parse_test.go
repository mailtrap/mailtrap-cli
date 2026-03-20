package send

import (
	"testing"
)

func TestParseEmailAddrPlain(t *testing.T) {
	addr, err := parseEmailAddr("email@test.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr.Email != "email@test.com" {
		t.Errorf("expected email 'email@test.com', got %q", addr.Email)
	}
	if addr.Name != "" {
		t.Errorf("expected empty name, got %q", addr.Name)
	}
}

func TestParseEmailAddrWithName(t *testing.T) {
	addr, err := parseEmailAddr("Name <email@test.com>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr.Email != "email@test.com" {
		t.Errorf("expected email 'email@test.com', got %q", addr.Email)
	}
	if addr.Name != "Name" {
		t.Errorf("expected name 'Name', got %q", addr.Name)
	}
}

func TestParseEmailAddrWithFullName(t *testing.T) {
	addr, err := parseEmailAddr("John Doe <john@example.com>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr.Email != "john@example.com" {
		t.Errorf("expected email 'john@example.com', got %q", addr.Email)
	}
	if addr.Name != "John Doe" {
		t.Errorf("expected name 'John Doe', got %q", addr.Name)
	}
}

func TestParseEmailAddrEmpty(t *testing.T) {
	_, err := parseEmailAddr("")
	if err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestParseEmailAddrWhitespace(t *testing.T) {
	addr, err := parseEmailAddr("  email@test.com  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr.Email != "email@test.com" {
		t.Errorf("expected trimmed email 'email@test.com', got %q", addr.Email)
	}
}

func TestParseEmailAddrs(t *testing.T) {
	addrs, err := parseEmailAddrs([]string{"a@test.com", "Name <b@test.com>"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(addrs) != 2 {
		t.Fatalf("expected 2 addresses, got %d", len(addrs))
	}
	if addrs[0].Email != "a@test.com" {
		t.Errorf("expected first email 'a@test.com', got %q", addrs[0].Email)
	}
	if addrs[1].Email != "b@test.com" {
		t.Errorf("expected second email 'b@test.com', got %q", addrs[1].Email)
	}
	if addrs[1].Name != "Name" {
		t.Errorf("expected second name 'Name', got %q", addrs[1].Name)
	}
}
