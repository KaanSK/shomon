package thehive

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestGetHiveClient(t *testing.T) {
	apiKey := "TeST"
	endpoint := "https://test.local"
	client, err := GetHiveClient(endpoint, apiKey, http.DefaultClient)
	if err != nil {
		t.Errorf(err.Error())
	}
	if client.apiKey != apiKey || client.url != endpoint {
		t.Errorf("Client could not be initalized with proper input")
	}
}

func TestCreateAlert(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	endpoint := "https://test.local"

	httpmock.RegisterResponder("POST", endpoint+"/api/v1/alert",
		func(req *http.Request) (*http.Response, error) {
			alert := Alert{}
			if err := json.NewDecoder(req.Body).Decode(&alert); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			resp, err := httpmock.NewJsonResponse(201, alert)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	hiveClient, err := GetHiveClient(endpoint, "TEST", http.DefaultClient)
	if err != nil {
		t.Errorf(err.Error())
	}
	alert := NewAlert()
	alert.Description = "TestDescription"
	alert.Type = "TEST_TYPE"
	alert.Title = "TestTitle"
	alert.Tags = []string{"test1", "test2"}
	alert.Source = "Shodan"
	alert.SourceRef = strconv.Itoa(100000 + rand.Intn(900000))
	alert.AddObservable("ip", "1.1.1.1")
	_, err = hiveClient.CreateAlert(alert)
	if err != nil {
		t.Errorf(err.Error())
	}

}
