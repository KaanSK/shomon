package producer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/shadowscatcher/shodan/models"
)

type ShodanWebhook struct {
	ShodanKey string
	Endpoint  string
	Port      int
}

func GetShodanWebhook(shodanKey string, endpoint string, port int) (*ShodanWebhook, error) {
	if shodanKey == "" {
		return nil, errors.New("empty API key")
	}

	if endpoint == "" {
		return nil, errors.New("endpoint is empty")
	}

	return &ShodanWebhook{
		ShodanKey: shodanKey,
		Endpoint:  endpoint,
		Port:      port,
	}, nil
}

func (sw *ShodanWebhook) ListenAlerts(ctx context.Context) (chan models.Service, error) {
	bannerChan := make(chan models.Service)

	bannerHandler := banner(bannerChan)
	http.HandleFunc(sw.Endpoint, bannerHandler)
	go func() {
		defer close(bannerChan)
		http.ListenAndServe(fmt.Sprintf(":%d", sw.Port), nil)
	}()

	return bannerChan, nil

}

func banner(bannerChan chan models.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		banner := models.Service{}
		/* alertID := req.Header.Get("SHODAN-ALERT-ID")
		alertName := req.Header.Get("SHODAN-ALERT-NAME")
		alertTrigger := req.Header.Get("SHODAN-ALERT-TRIGGER")
		alertVerify := req.Header.Get("SHODAN-SIGNATURE-SHA1")
		if containsEmpty(alertID, alertName, alertTrigger, alertVerify) {
			w.WriteHeader(http.StatusBadRequest)
			return
		} */
		if err := json.NewDecoder(req.Body).Decode(&banner); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			bannerChan <- banner
			w.WriteHeader(http.StatusOK)
		}
	}
}

func containsEmpty(ss ...string) bool {
	for _, s := range ss {
		if s == "" {
			return true
		}
	}
	return false
}
