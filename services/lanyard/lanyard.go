package lanyard

import (
	"context"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
)

type Service struct {
	http      *req.Client
	defaultID string
}

func (s *Service) Get(ctx context.Context, id string) (*Presence, error) {
	var res GetPresenceResponse

	_, err := s.http.R().
		SetContext(ctx).
		SetSuccessResult(&res).
		Get("/users/" + id)

	if err != nil {
		return nil, fmt.Errorf("failed to get presence: %w", err)
	}

	if !res.Success {
		return nil, fmt.Errorf("lanyard internal error")
	}

	p := &res.Data

	return p, nil
}

func (s *Service) GetDefault(ctx context.Context) (*Presence, error) {
	return s.Get(ctx, s.defaultID)
}

type GetPresenceResponse struct {
	Data    Presence `json:"data"`
	Success bool     `json:"success,omitempty"`
}

type Presence struct {
	User                    User       `json:"discord_user"`
	Activities              []Activity `json:"activities,omitempty"`
	DiscordStatus           string     `json:"discord_status,omitempty"`
	ActiveOnDiscordWeb      bool       `json:"active_on_discord_web,omitempty"`
	ActiveOnDiscordDesktop  bool       `json:"active_on_discord_desktop,omitempty"`
	ActiveOnDiscordMobile   bool       `json:"active_on_discord_mobile,omitempty"`
	ActiveOnDiscordEmbedded bool       `json:"active_on_discord_embedded,omitempty"`
	ActiveOnDiscordVr       bool       `json:"active_on_discord_vr,omitempty"`
	ListeningToSpotify      bool       `json:"listening_to_spotify,omitempty"`
	Spotify                 any        `json:"spotify,omitempty"`
}

type User struct {
	Avatar               string `json:"avatar,omitempty"`
	AvatarDecorationData any    `json:"avatar_decoration_data,omitempty"`
	Bot                  bool   `json:"bot,omitempty"`
	Collectibles         any    `json:"collectibles,omitempty"`
	Discriminator        string `json:"discriminator,omitempty"`
	DisplayName          string `json:"display_name,omitempty"`
	DisplayNameStyles    any    `json:"display_name_styles,omitempty"`
	GlobalName           string `json:"global_name,omitempty"`
	ID                   string `json:"id,omitempty"`
	PrimaryGuild         any    `json:"primary_guild,omitempty"`
	PublicFlags          int    `json:"public_flags,omitempty"`
	Username             string `json:"username,omitempty"`
}

type Activity struct {
	ApplicationID string `json:"application_id,omitempty"`
	Assets        struct {
		LargeImage string `json:"large_image,omitempty"`
		LargeText  string `json:"large_text,omitempty"`
		SmallImage string `json:"small_image,omitempty"`
		SmallText  string `json:"small_text,omitempty"`
	} `json:"assets"`
	CreatedAt  int64  `json:"created_at,omitempty"`
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Details    string `json:"details,omitempty"`
	Platform   string `json:"platform,omitempty"`
	SessionID  string `json:"session_id,omitempty"`
	State      string `json:"state,omitempty"`
	Timestamps struct {
		Start int64 `json:"start,omitempty"`
	} `json:"timestamps"`
	Type int `json:"type,omitempty"`
}

func WithDefaultID(id string) Option {
	return func(s *Service) {
		s.defaultID = id
	}
}

func WithBaseURL(url string) Option {
	return func(s *Service) {
		s.http.SetBaseURL(url)
	}
}

type Option func(*Service)

func New(opts ...Option) *Service {
	s := &Service{
		http: req.NewClient().
			SetBaseURL("https://api.lanyard.rest/v1/").
			SetTimeout(time.Second),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}
