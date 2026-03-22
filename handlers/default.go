package handlers

import (
	views "git.jbennett.dev/persona-www/components"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Option func(*DefaultHandler)

func WithLogger(logger zerolog.Logger) Option {
	return func(h *DefaultHandler) {
		h.logger = logger
	}
}

type DefaultHandler struct {
	logger    zerolog.Logger
	lanyardID string
}

func New(lanyardID string, opts ...Option) (*DefaultHandler, error) {
	h := &DefaultHandler{
		logger:    log.Logger,
		lanyardID: lanyardID,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h, nil
}

func (h *DefaultHandler) Get(c echo.Context) error {
	return render(c, views.Index("jabenne.net", h.lanyardID))
}

func (h *DefaultHandler) GetPresence(c echo.Context) error {
	return render(c, views.Presence(h.lanyardID))
}

func render(ctx echo.Context, cmp templ.Component) error {
	return cmp.Render(ctx.Request().Context(), ctx.Response())
}
