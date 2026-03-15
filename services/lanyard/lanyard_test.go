package lanyard_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.jbennett.dev/persona-www/services/lanyard"
)

func newTestServer(t *testing.T, handler http.Handler) (*httptest.Server, *lanyard.Service) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	svc := lanyard.New(
		lanyard.WithBaseURL(srv.URL+"/"),
		lanyard.WithDefaultID("testuser"),
	)
	return srv, svc
}

func TestService_Get_Success(t *testing.T) {
	want := lanyard.Presence{
		DiscordStatus: "online",
		Activities: []lanyard.Activity{
			{Name: "Visual Studio Code", Details: "editing main.go", State: "In a folder"},
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/testuser" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		resp := map[string]any{
			"success": true,
			"data":    want,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	_, svc := newTestServer(t, handler)
	got, err := svc.Get(context.Background(), "testuser")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.DiscordStatus != want.DiscordStatus {
		t.Errorf("DiscordStatus = %q, want %q", got.DiscordStatus, want.DiscordStatus)
	}
	if len(got.Activities) != 1 {
		t.Fatalf("len(Activities) = %d, want 1", len(got.Activities))
	}
	if got.Activities[0].Name != want.Activities[0].Name {
		t.Errorf("Activity.Name = %q, want %q", got.Activities[0].Name, want.Activities[0].Name)
	}
	if got.Activities[0].Details != want.Activities[0].Details {
		t.Errorf("Activity.Details = %q, want %q", got.Activities[0].Details, want.Activities[0].Details)
	}
}

func TestService_Get_RequestError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close() // close immediately so requests fail

	svc := lanyard.New(
		lanyard.WithBaseURL(srv.URL+"/"),
		lanyard.WithDefaultID("testuser"),
	)

	_, err := svc.Get(context.Background(), "testuser")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to get presence") {
		t.Errorf("error message %q should contain \"failed to get presence\"", err.Error())
	}
}

func TestService_Get_LanyardFailure(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{"success": false}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	_, svc := newTestServer(t, handler)
	_, err := svc.Get(context.Background(), "testuser")
	if err == nil {
		t.Fatal("expected error for success=false, got nil")
	}
	if !strings.Contains(err.Error(), "lanyard internal error") {
		t.Errorf("error message %q should contain \"lanyard internal error\"", err.Error())
	}
}

func TestService_GetDefault_UsesConfiguredID(t *testing.T) {
	const wantID = "myspecialid"
	var gotPath string

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		resp := map[string]any{
			"success": true,
			"data":    lanyard.Presence{DiscordStatus: "online"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	svc := lanyard.New(
		lanyard.WithBaseURL(srv.URL+"/"),
		lanyard.WithDefaultID(wantID),
	)

	_, err := svc.GetDefault(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantPath := "/users/" + wantID
	if gotPath != wantPath {
		t.Errorf("request path = %q, want %q", gotPath, wantPath)
	}
}

func TestService_WithBaseURL_ConfiguresClient(t *testing.T) {
	reached := false

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached = true
		resp := map[string]any{
			"success": true,
			"data":    lanyard.Presence{DiscordStatus: "idle"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	_, svc := newTestServer(t, handler)
	_, err := svc.GetDefault(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reached {
		t.Error("WithBaseURL did not configure the HTTP client to reach the test server")
	}
}
