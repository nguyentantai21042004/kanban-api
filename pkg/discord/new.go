package discord

import (
	"errors"

	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
)

type DiscordWebhook struct {
	ID    string
	Token string
}

func NewDiscordWebhook(id, token string) (*DiscordWebhook, error) {
	if id == "" || token == "" {
		return nil, errors.New("id and token are required")
	}

	return &DiscordWebhook{
		ID:    id,
		Token: token,
	}, nil
}

type Discord struct {
	l       log.Logger
	webhook *DiscordWebhook
}

func New(l log.Logger, webhook *DiscordWebhook) (*Discord, error) {
	if webhook == nil {
		return nil, errors.New("webhook is required")
	}

	return &Discord{
		l:       l,
		webhook: webhook,
	}, nil
}
