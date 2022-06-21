package thehive

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type TheHiveClient struct {
	apiKey string
	url    string
	HTTP   *http.Client
	Logger log.Logger
}

func GetHiveClient(url string, key string, client *http.Client) (*TheHiveClient, error) {
	/* 	if key == "" {
	   		return nil, errors.New("empty Hive API key")
	   	}

	   	if client == nil {
	   		return nil, errors.New("HTTP client is nil")
	   	} */

	return &TheHiveClient{
		apiKey: key,
		url:    url,
		HTTP:   client,
	}, nil
}

func NewAlert() Alert {
	return Alert{}
}

func (a *Alert) AddObservable(obsType string, obs string) {
	obsInstance := Observable{
		Data:     obs,
		DataType: obsType,
	}
	a.Observables = append(a.Observables, obsInstance)
}

func (s *TheHiveClient) CreateAlert(alert Alert) (id string, err error) {
	if s == nil {
		return id, errors.New("not initialized hive client")
	}
	payload, err := json.Marshal(alert)
	if err != nil {
		return id, err
	}
	req, err := http.NewRequest("POST", s.url+"/api/v1/alert", bytes.NewBuffer(payload))
	if err != nil {
		return id, err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTP.Do(req)
	if err != nil {
		return id, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		if err != nil {
			return id, err
		}
		return id, errors.New(string(body))
	}
	createdAlert := NewAlert()
	err = json.Unmarshal(body, &createdAlert)
	if err != nil {
		return id, errors.New(err.Error())
	}
	return createdAlert.Id, nil
}

type DeleteAlertsInput struct {
	Ids []string `json:"ids"`
}

func (s *TheHiveClient) DeleteAlerts(ids []string) error {
	if s == nil {
		return errors.New("not initialized hive client")
	}
	payload, err := json.Marshal(&DeleteAlertsInput{Ids: ids})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.url+"/api/v1/alert/delete/_bulk", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	return nil
}
