package handlers

import (
	"context"

	views "git.jbennett.dev/persona-www/components"
	"git.jbennett.dev/persona-www/services/lanyard"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LanyardService interface {
	GetDefault(context.Context) (*lanyard.Presence, error)
}

type Option func(*DefaultHandler)

func WithLogger(logger zerolog.Logger) Option {
	return func(h *DefaultHandler) {
		h.logger = logger
	}
}

type DefaultHandler struct {
	logger  zerolog.Logger
	lanyard LanyardService
}

func New(lanyard LanyardService, opts ...Option) (*DefaultHandler, error) {
	h := &DefaultHandler{
		logger:  log.Logger,
		lanyard: lanyard,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h, nil
}

func (h *DefaultHandler) Get(c echo.Context) error {
	return render(c, views.Index("jabenne.net"))
}

func (h *DefaultHandler) GetPresence(c echo.Context) error {
	p, err := h.lanyard.GetDefault(c.Request().Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to fetch lanyard presence")
		return render(c, views.Presence("offline", nil))
	}
	var activity *lanyard.Activity
	if len(p.Activities) > 0 {
		activity = &p.Activities[0]
	}

	return render(c, views.Presence(p.DiscordStatus, activity))
}

func render(ctx echo.Context, cmp templ.Component) error {
	return cmp.Render(ctx.Request().Context(), ctx.Response())
}
