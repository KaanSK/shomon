package service

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/KaanSK/shomon/log"
	"github.com/KaanSK/shomon/thehive"
	"github.com/jarcoal/httpmock"
	"github.com/shadowscatcher/shodan/models"
)

type MockStream struct {
}

func (ss *MockStream) ListenAlerts(ctx context.Context) (chan models.Service, error) {
	bannerChan := make(chan models.Service)
	bannerJson := `{
		"hash": 1015805840,
		"timestamp": "2021-01-28T04:16:08.387364",
		"hostnames": [
			"177-70-193-184-msltr-cw-1.visaonet.com.br"
		],
		"org": "TESTORG",
		"data": "SIP/2.0 404 Not Found\r\nFrom: ;tag=root\r\nTo: ;tag=b235f0-b146c1b8-13c4-50029-ec2e4-6c44dfcf-ec2e4\r\nCall-ID: 50000\r\nCSeq: 42 OPTIONS\r\nVia: SIP/2.0/UDP nm;received=224.238.62.40;rport=26810;branch=foo\r\nSupported: replaces,100rel,timer\r\nAccept: application/sdp\r\nAllow: INVITE,ACK,CANCEL,BYE,OPTIONS,REFER,INFO,NOTIFY,PRACK,MESSAGE\r\nContent-Length: 0\r\n\r\n",
		"port": 5060,
		"transport": "udp",
		"info": "SIP end point; Status: 404 Not Found",
		"isp": "L M Tiko Kamide - Sva",
		"asn": "AS28359",
		"location": {
			"country_code3": null,
			"city": "Jardim Alegre",
			"region_code": "PR",
			"postal_code": null,
			"longitude": -51.7213,
			"country_code": "BR",
			"latitude": -24.2123,
			"country_name": "Brazil",
			"area_code": null,
			"dma_code": null
		},
		"ip": 2974204344,
		"domains": [
			"visaonet.com.br"
		],
		"ip_str": "177.70.193.184",
		"_id": "45ad6383-1b1d-4c5d-8584-d586fbdefbc3",
		"os": null,
		"_shodan": {
			"crawler": "bf213bc419cc8491376c12af31e32623c1b6f467",
			"options": {},
			"id": "220ef463-756f-4446-a89f-685053da8865",
			"module": "sip",
			"ptr": true
		},
		"opts": {}
	}`
	banner := models.Service{}
	buf := make([]byte, 4)
	err := json.Unmarshal([]byte(bannerJson), &banner)
	if err != nil {
		return nil, err
	}
	rand.Seed(time.Now().UTC().UnixNano())
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(1 * time.Second)
			ip := rand.Uint32()
			binary.LittleEndian.PutUint32(buf, ip)
			banner.IPstr = fmt.Sprintf("%s", net.IP(buf))
			bannerChan <- banner
		}
		close(bannerChan)
	}()

	return bannerChan, nil
}

func TestListenStreamAlerts(t *testing.T) {
	ms := &MockStream{}
	var wg sync.WaitGroup
	ctx := context.Background()
	logger, err := log.New(os.Stdout, "ERROR")
	if err != nil {
		log.Fatal(err.Error())
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	endpoint := "https://test.local"

	httpmock.RegisterResponder("POST", endpoint+"/api/v1/alert",
		func(req *http.Request) (*http.Response, error) {
			alert := thehive.Alert{}
			if err := json.NewDecoder(req.Body).Decode(&alert); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			alert.Id = fmt.Sprintf("~%d", rand.Intn(1000000))
			resp, err := httpmock.NewJsonResponse(201, alert)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	hiveClient, err := thehive.GetHiveClient(endpoint, "TEST", http.DefaultClient)
	if err != nil {
		t.Errorf(err.Error())
	}

	srv := Service{
		ShodanClient: ms,
		HiveClient:   hiveClient,
		wg:           &wg,
		ctx:          ctx,
		Logger:       *logger,
	}

	alertChan, err := srv.ShodanClient.ListenAlerts(ctx)
	if err != nil {
		t.Errorf(err.Error())
	}

	for banner := range alertChan {
		srv.wg.Add(1)
		go func() {
			_, err := srv.ProcessAlert(banner)
			if err != nil {
				t.Errorf(err.Error())
			}
		}()
	}
	wg.Wait()
}
