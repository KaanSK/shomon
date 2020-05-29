package shodan

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	conf "github.com/KaanSK/shomon/pkg/conf"
	lw "github.com/KaanSK/shomon/pkg/logwrapper"
	"github.com/KaanSK/shomon/pkg/thehive"
)

var (
	errShomonServiceStop = errors.New("listener service stopped")
	logger               = lw.NewLogger()
)

func parseResponse(destination interface{}, body io.Reader) error {
	var err error

	if w, ok := destination.(io.Writer); ok {
		_, err = io.Copy(w, body)
	} else {
		decoder := json.NewDecoder(body)
		err = decoder.Decode(destination)
	}

	return err
}

func handleAlertStream(ch chan *HostData) {
	defer func() {
		close(ch)
	}()
	resp, err := http.Get("https://stream.shodan.io/shodan/alert?key=" + conf.Config.ShodanKey)
	if err != nil {
		logger.Error(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = GetErrorFromResponse(resp)
		resp.Body.Close()
		logger.Fatal(err)
	}

	reader := bufio.NewReader(resp.Body)
	for {
		banner := new(HostData)
		chunk, err := reader.ReadBytes('\n')
		if err != nil {
			resp.Body.Close()
			break
		}

		chunk = bytes.TrimRight(chunk, "\n\r")
		if len(chunk) == 0 {
			continue
		}

		if err := parseResponse(banner, bytes.NewBuffer(chunk)); err != nil {
			resp.Body.Close()
			logger.Error(err)
			break
		}

		ch <- banner
	}
}

// ListenAlerts : Used to listen streaming monitoring API
func ListenAlerts() error {
	ch := make(chan *HostData)
	go handleAlertStream(ch)

	logger.Info("listening process initiated")

	for {
		banner, ok := <-ch
		if !ok {
			break
		}

		hiveAlert := new(thehive.HiveAlert)
		foundService := fmt.Sprintf("%s:%d", banner.IP, banner.Port)
		hiveAlert.Title = fmt.Sprintf("Alert: %s", foundService)
		hiveAlert.Description = "Test description"

		if conf.Config.CaseTemplate != "" {
			hiveAlert.CaseTemplate = conf.Config.CaseTemplate
		}

		hiveAlert.Source = "Shodan"
		hash := md5.Sum([]byte(foundService))
		hiveAlert.SourceRef = hex.EncodeToString(hash[:])

		logger.Info("triggered alarm for: " + hiveAlert.SourceRef)

		err := thehive.CreateAlert(hiveAlert)
		if err != nil {
			return err
		}
		logger.Info("created alert for " + hiveAlert.SourceRef)
	}
	return errShomonServiceStop
}
