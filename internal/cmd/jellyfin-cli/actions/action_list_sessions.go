package actions

import (
	"context"
	"log"
	"time"

	"github.com/dustin/go-humanize"

	"codeberg.org/jfenske/jellyfin-cli/api"
)

type listSessionsExecutorImpl struct {
	client api.JellyfinApiClient
}

func NewListSessionsExecutor(client api.JellyfinApiClient) Executor {
	return &listSessionsExecutorImpl{
		client: client,
	}
}

func (e *listSessionsExecutorImpl) Run(ctx context.Context, arguments map[string]interface{}) error {
	if sessions, err := e.client.ListSessions(ctx); err != nil {
		return err
	} else if len(sessions) == 0 {
		log.Print("No sessions")
	} else {
		var active []api.Session
		for _, session := range sessions {
			if time.Since(session.Date) <= 10*time.Minute {
				active = append(active, session)
			}
		}

		if len(active) == 0 {
			log.Print("No sessions")
		} else {
			log.Print("Sessions:")
			for _, session := range active {
				duration := humanize.RelTime(time.Now(), session.Date, "", "ago")
				log.Printf(" - %s (%s) %s", session.User, session.Device, duration)
			}
		}
	}

	return nil
}
