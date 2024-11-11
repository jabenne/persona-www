package handlers

import (
	"math/rand/v2"
	"time"

	views "git.jbennett.dev/persona-www/components"
	"git.jbennett.dev/persona-www/services"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type GitService interface {
    GetHistorySince(string, time.Time) ([]int, error) 
}

type DefaultHandler struct {
    name string
    git GitService 
}

func New() (*DefaultHandler, error) {
    ghC, err := services.NewGithubConfigFromEnv()
    if err != nil {
        panic(err)
    }
    return &DefaultHandler{
        git: services.NewGithubService(ghC),
    }, nil
}

func (h *DefaultHandler) Get(c echo.Context) error {
    gh, err := h.git.GetHistorySince("jabenne", time.Now().AddDate(0, 0, -364))
    if err != nil {
        return err
    }
    return render(c, views.Index("Test", gh))
}

func render(ctx echo.Context, cmp templ.Component) error {
    return cmp.Render(ctx.Request().Context(), ctx.Response())
}

func fudgeGitHistory() []int {
    gitHistory := make([]int, 364)
    for i := range gitHistory {
        gitHistory[i] = rand.IntN(3)
    }
    return gitHistory
}
