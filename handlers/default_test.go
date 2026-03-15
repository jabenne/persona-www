package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.jbennett.dev/persona-www/handlers"
	"git.jbennett.dev/persona-www/services/lanyard"
	"github.com/labstack/echo/v4"
)

// fakeLanyard is a test double for LanyardService.
type fakeLanyard struct {
	presence    *lanyard.Presence
	err         error
	capturedCtx context.Context
}

func (f *fakeLanyard) GetDefault(ctx context.Context) (*lanyard.Presence, error) {
	f.capturedCtx = ctx
	return f.presence, f.err
}

// newEchoCtx creates a minimal Echo context for GET requests.
func newEchoCtx(t *testing.T) (echo.Context, *httptest.ResponseRecorder) {
	t.Helper()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func newHandler(t *testing.T, svc *fakeLanyard) *handlers.DefaultHandler {
	t.Helper()
	h, err := handlers.New(svc)
	if err != nil {
		t.Fatalf("handlers.New: %v", err)
	}
	return h
}

func TestDefaultHandler_Get_RendersIndex(t *testing.T) {
	svc := &fakeLanyard{}
	h := newHandler(t, svc)
	c, rec := newEchoCtx(t)

	if err := h.Get(c); err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	body := rec.Body.String()
	for _, want := range []string{"jabenne.net", "Jack Bennett", "jabenne.net"} {
		if !strings.Contains(body, want) {
			t.Errorf("response body missing %q", want)
		}
	}
}

func TestDefaultHandler_GetPresence_SuccessWithActivity(t *testing.T) {
	svc := &fakeLanyard{
		presence: &lanyard.Presence{
			DiscordStatus: "online",
			Activities: []lanyard.Activity{
				{Name: "Visual Studio Code", Details: "editing main.go", State: "In a folder"},
			},
		},
	}
	h := newHandler(t, svc)
	c, rec := newEchoCtx(t)

	if err := h.GetPresence(c); err != nil {
		t.Fatalf("GetPresence returned error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Online") {
		t.Error("response body missing \"Online\"")
	}
	if !strings.Contains(body, "Visual Studio Code") {
		t.Error("response body missing activity name \"Visual Studio Code\"")
	}
	if !strings.Contains(body, "editing main.go") {
		t.Error("response body missing activity details \"editing main.go\"")
	}
}

func TestDefaultHandler_GetPresence_SuccessNoActivity(t *testing.T) {
	svc := &fakeLanyard{
		presence: &lanyard.Presence{
			DiscordStatus: "idle",
			Activities:    []lanyard.Activity{},
		},
	}
	h := newHandler(t, svc)
	c, rec := newEchoCtx(t)

	if err := h.GetPresence(c); err != nil {
		t.Fatalf("GetPresence returned error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Idle") {
		t.Error("response body missing \"Idle\"")
	}
	// No activity detail spans should be present.
	if strings.Contains(body, "Visual Studio Code") {
		t.Error("response body should not contain activity name when no activities")
	}
}

func TestDefaultHandler_GetPresence_ServiceErrorFallsBackOffline(t *testing.T) {
	svc := &fakeLanyard{
		err: errors.New("connection refused"),
	}
	h := newHandler(t, svc)
	c, rec := newEchoCtx(t)

	if err := h.GetPresence(c); err != nil {
		t.Fatalf("GetPresence returned error (expected graceful fallback): %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Offline") {
		t.Error("response body missing \"Offline\" for error fallback")
	}
}

func TestDefaultHandler_GetPresence_UsesRequestContext(t *testing.T) {
	svc := &fakeLanyard{
		presence: &lanyard.Presence{DiscordStatus: "online"},
	}
	h := newHandler(t, svc)
	c, _ := newEchoCtx(t)

	if err := h.GetPresence(c); err != nil {
		t.Fatalf("GetPresence returned error: %v", err)
	}
	if svc.capturedCtx == nil {
		t.Error("GetDefault was not called with a context")
	}
}
