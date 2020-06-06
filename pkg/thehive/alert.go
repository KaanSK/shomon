package thehive

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	conf "github.com/KaanSK/shomon/pkg/conf"
)

// HiveAlert : Alert structure
type HiveAlert struct {
	CaseTemplate string        `json:"caseTemplate,omitempty"`
	Artifacts    []interface{} `json:"artifacts"`
	CreatedAt    int64         `json:"createdAt"`
	CreatedBy    string        `json:"createdBy"`
	Date         int64         `json:"date"`
	Description  string        `json:"description"`
	Follow       bool          `json:"follow"`
	ID           string        `json:"id,omitempty"`
	LastSyncDate int64         `json:"lastSyncDate"`
	Severity     int           `json:"severity"`
	Source       string        `json:"source"`
	SourceRef    string        `json:"sourceRef"`
	Status       string        `json:"status"`
	Title        string        `json:"title"`
	Tlp          int           `json:"tlp"`
	Type         string        `json:"type"`
	User         string        `json:"user,omitempty"`
}

func getResponseJSON(url string, target interface{}, input interface{}) (interface{}, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	payload, err := json.Marshal(input)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+conf.Config.TheHiveKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := netClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		bodyString := string(bodyBytes)
		return nil, errors.New(bodyString)
	}

	err = json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		return nil, err
	}
	return target, nil

}

// CreateAlert : Used to create alert on thehive
func CreateAlert(alertObject *HiveAlert) error {
	endpoint := conf.Config.Endpoint
	_, err := getResponseJSON(endpoint, HiveAlert{}, *alertObject)
	if err != nil {
		return err
	}
	return nil
}
