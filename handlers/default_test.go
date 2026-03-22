package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.jbennett.dev/persona-www/handlers"
	"github.com/labstack/echo/v4"
)

// newEchoCtx creates a minimal Echo context for GET requests.
func newEchoCtx(t *testing.T) (echo.Context, *httptest.ResponseRecorder) {
	t.Helper()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func newHandler(t *testing.T, lanyardID string) *handlers.DefaultHandler {
	t.Helper()
	h, err := handlers.New(lanyardID)
	if err != nil {
		t.Fatalf("handlers.New: %v", err)
	}
	return h
}

func TestDefaultHandler_Get_RendersIndex(t *testing.T) {
	h := newHandler(t, "")
	c, rec := newEchoCtx(t)

	if err := h.Get(c); err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	body := rec.Body.String()
	for _, want := range []string{"jabenne.net", "Jack Bennett"} {
		if !strings.Contains(body, want) {
			t.Errorf("response body missing %q", want)
		}
	}
}

func TestDefaultHandler_GetPresence_RendersWithDiscordID(t *testing.T) {
	const discordID = "testuser123"
	h := newHandler(t, discordID)
	c, rec := newEchoCtx(t)

	if err := h.GetPresence(c); err != nil {
		t.Fatalf("GetPresence returned error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	body := rec.Body.String()
	if !strings.Contains(body, discordID) {
		t.Errorf("response body missing discord ID %q", discordID)
	}
}

func TestDefaultHandler_GetPresence_EmptyIDRendersPresence(t *testing.T) {
	h := newHandler(t, "")
	c, rec := newEchoCtx(t)

	if err := h.GetPresence(c); err != nil {
		t.Fatalf("GetPresence returned error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}
