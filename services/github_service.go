package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v66/github"
)

type GithubService struct {
    client *github.Client
}

type GithubConfig struct {
    Token string
}

func NewGithubConfigFromEnv() (*GithubConfig, error) {
    pat, exists := os.LookupEnv("GITHUB_PAT")
    if !exists {
        return nil, fmt.Errorf("GITHUB_PAT not set")
    }

    return &GithubConfig{
        Token: pat,
    }, nil
}

func NewGithubService(c *GithubConfig) *GithubService{
    client := github.NewClient(nil).WithAuthToken(c.Token)
    return &GithubService{
        client: client,
    }
}

func (s *GithubService) GetHistorySince(user string, t time.Time) ([]int, error) {
    page := 1
    records := 100
    tUnix := t.Unix()
    now := time.Now()
    
    var history []*github.Event
    var numHistory [364]int
    out:
    for {
        opts := github.ListOptions{
            Page: page,
            PerPage: records,
        }
        events, _, err := s.client.Activity.ListEventsPerformedByUser(context.Background(), user, false, &opts)
        if err != nil {
            return nil, err
        }
        if len(events) == 0 {
            break out
        }
        for _, e := range events {
            if e.CreatedAt.Unix() < tUnix {
                break out
            } else {
                history = append(history, e) 
                diff := now.Sub(e.CreatedAt.Time)
                index := int(diff.Hours() / 24)
                numHistory[index] = numHistory[index] + 1
            }
        }
        page = page + 1
    }

    return numHistory[:], nil
}
