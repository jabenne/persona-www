package services_test

import (
	"fmt"
	"testing"
	"time"

	"git.jbennett.dev/persona-www/services"
)

func TestGetAnnualHistory(t *testing.T) {
    c, err:= services.NewGithubConfigFromEnv()
    if err != nil {
        t.Fatal(err)
    }
    s := services.NewGithubService(c)
    r, err := s.GetHistorySince("jabenne", time.Now().AddDate(-1, 0, 0))
    if err != nil {
        t.Fatal(err)
    }
    fmt.Println(r)
}
