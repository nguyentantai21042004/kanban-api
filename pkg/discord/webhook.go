package discord

import (
	"context"
	"fmt"

	"gitlab.com/tantai-kanban/kanban-api/pkg/curl"
)

func (d *Discord) ReportBug(ctx context.Context, message string) error {
	url := fmt.Sprintf(webhookURL, d.webhook.ID, d.webhook.Token)

	header := map[string]string{
		"Content-Type": "application/json",
	}

	jsonData := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title":       "SMAP Service Error Report",
				"description": fmt.Sprintf("```%s```", message),
				"color":       15158332, // Red color for error
			},
		},
	}

	// Send the jsonData directly without marshaling it first
	_, err := curl.Post(url, header, jsonData)
	if err != nil {
		d.l.Errorf(ctx, "pkg.discord.webhook.ReportBug.curl.Post: %v", err)
		return err
	}

	return nil
}
