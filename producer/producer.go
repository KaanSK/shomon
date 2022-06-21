package producer

import (
	"context"

	"github.com/shadowscatcher/shodan/models"
)

type Producer interface {
	ListenAlerts(ctx context.Context) (chan models.Service, error)
}
