package stats_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/commands/stats"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/spf13/viper"
)

func setupTest(handler http.HandlerFunc) (*cmdutil.Factory, *bytes.Buffer, func()) {
	server := httptest.NewServer(handler)

	c := client.New("test-token")
	c.SetBaseURL(client.BaseGeneral, server.URL)

	buf := &bytes.Buffer{}
	f := &cmdutil.Factory{
		Config: func() *config.Config {
			return &config.Config{APIToken: "test-token", AccountID: "123"}
		},
		IOStreams: &cmdutil.IOStreams{
			Out:    buf,
			ErrOut: &bytes.Buffer{},
		},
		ClientOverride: c,
	}

	viper.Set("api-token", "test-token")
	viper.Set("account-id", "123")
	viper.Set("output", "table")

	return f, buf, func() {
		server.Close()
		viper.Reset()
	}
}

func TestStatsGet(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "stats") {
			t.Errorf("expected path to contain 'stats', got %s", r.URL.Path)
		}
		if r.URL.Query().Get("start_date") != "2024-01-01" {
			t.Errorf("expected start_date=2024-01-01, got %s", r.URL.Query().Get("start_date"))
		}
		if r.URL.Query().Get("end_date") != "2024-01-31" {
			t.Errorf("expected end_date=2024-01-31, got %s", r.URL.Query().Get("end_date"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"delivery_count": 100,
			"delivery_rate":  0.95,
			"bounce_count":   5,
			"bounce_rate":    0.05,
			"open_count":     50,
			"open_rate":      0.5,
			"click_count":    20,
			"click_rate":     0.2,
			"spam_count":     1,
			"spam_rate":      0.01,
		})
	})
	defer cleanup()

	cmd := stats.NewCmdStats(f)
	cmd.SetArgs([]string{"get", "--start-date", "2024-01-01", "--end-date", "2024-01-31"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "100") {
		t.Errorf("expected output to contain '100', got:\n%s", output)
	}
	if !strings.Contains(output, "0.95") {
		t.Errorf("expected output to contain '0.95', got:\n%s", output)
	}
}

func TestStatsGetMissingDates(t *testing.T) {
	f, _, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {})
	defer cleanup()

	buf := &bytes.Buffer{}
	f.IOStreams.Out = buf

	cmd := stats.NewCmdStats(f)
	cmd.SetArgs([]string{"get"})
	cmd.SetOut(buf)
	cmd.SetErr(&bytes.Buffer{})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when required flags are missing")
	}
}

func TestStatsByDomain(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "stats") || !strings.Contains(r.URL.Path, "domains") {
			t.Errorf("expected path to contain 'stats' and 'domains', got %s", r.URL.Path)
		}
		if r.URL.Query().Get("start_date") != "2024-01-01" {
			t.Errorf("expected start_date=2024-01-01, got %s", r.URL.Query().Get("start_date"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"delivery_count": 50, "delivery_rate": 0.9, "bounce_count": 2, "bounce_rate": 0.04, "open_count": 25, "open_rate": 0.5, "click_count": 10, "click_rate": 0.2, "spam_count": 0, "spam_rate": 0.0},
		})
	})
	defer cleanup()

	cmd := stats.NewCmdStats(f)
	cmd.SetArgs([]string{"by-domain", "--start-date", "2024-01-01", "--end-date", "2024-01-31"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "50") {
		t.Errorf("expected output to contain '50', got:\n%s", output)
	}
}

func TestStatsByCategory(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "stats") || !strings.Contains(r.URL.Path, "categories") {
			t.Errorf("expected path to contain 'stats' and 'categories', got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"delivery_count": 30, "delivery_rate": 0.85, "bounce_count": 3, "bounce_rate": 0.06, "open_count": 15, "open_rate": 0.5, "click_count": 8, "click_rate": 0.27, "spam_count": 0, "spam_rate": 0.0},
		})
	})
	defer cleanup()

	cmd := stats.NewCmdStats(f)
	cmd.SetArgs([]string{"by-category", "--start-date", "2024-01-01", "--end-date", "2024-01-31"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "30") {
		t.Errorf("expected output to contain '30', got:\n%s", output)
	}
}

func TestStatsByEsp(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "stats") || !strings.Contains(r.URL.Path, "email_service_providers") {
			t.Errorf("expected path to contain 'stats' and 'email_service_providers', got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"delivery_count": 40, "delivery_rate": 0.88, "bounce_count": 4, "bounce_rate": 0.08, "open_count": 20, "open_rate": 0.5, "click_count": 12, "click_rate": 0.3, "spam_count": 1, "spam_rate": 0.02},
		})
	})
	defer cleanup()

	cmd := stats.NewCmdStats(f)
	cmd.SetArgs([]string{"by-esp", "--start-date", "2024-01-01", "--end-date", "2024-01-31"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "40") {
		t.Errorf("expected output to contain '40', got:\n%s", output)
	}
}

func TestStatsByDate(t *testing.T) {
	f, buf, cleanup := setupTest(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "stats") || !strings.Contains(r.URL.Path, "date") {
			t.Errorf("expected path to contain 'stats' and 'date', got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"delivery_count": 60, "delivery_rate": 0.92, "bounce_count": 6, "bounce_rate": 0.06, "open_count": 30, "open_rate": 0.5, "click_count": 15, "click_rate": 0.25, "spam_count": 2, "spam_rate": 0.03},
		})
	})
	defer cleanup()

	cmd := stats.NewCmdStats(f)
	cmd.SetArgs([]string{"by-date", "--start-date", "2024-01-01", "--end-date", "2024-01-31"})
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "60") {
		t.Errorf("expected output to contain '60', got:\n%s", output)
	}
}
