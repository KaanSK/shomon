package producer

import (
	"context"
	"errors"
	"net/http"

	"github.com/shadowscatcher/shodan"
	"github.com/shadowscatcher/shodan/models"
)

type ShodanStream struct {
	client *shodan.StreamClient
}

func GetShodanStreamClient(ShodanKey string, httpClient *http.Client) (*ShodanStream, error) {
	if ShodanKey == "" {
		return nil, errors.New("empty Shodan API key")
	}

	if httpClient == nil {
		return nil, errors.New("HTTP client is nil")
	}

	client, err := shodan.GetStreamClient(ShodanKey, httpClient)
	if err != nil {
		return nil, err
	}

	return &ShodanStream{
		client: client,
	}, nil
}

func (ss *ShodanStream) ListenAlerts(ctx context.Context) (chan models.Service, error) {
	alertChan, err := ss.client.Alerts(ctx)
	if err != nil {
		return nil, err
	}
	return alertChan, nil
}
