package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/KaanSK/shomon/conf"
	"github.com/KaanSK/shomon/log"
	"github.com/KaanSK/shomon/producer"
	"github.com/KaanSK/shomon/thehive"
	"github.com/shadowscatcher/shodan/models"
)

type Result struct {
	Message string
	Error   error
}

type Service struct {
	Config       conf.ShomonConfig
	HiveClient   *thehive.TheHiveClient
	ShodanClient producer.Producer
	Logger       log.Logger
	wg           *sync.WaitGroup
	ctx          context.Context
}

func New(wg *sync.WaitGroup, ctx context.Context) (srv Service, err error) {
	shomonConf, err := conf.New()
	if err != nil {
		log.Fatal(err.Error())
	}
	logger, err := log.New(os.Stdout, shomonConf.LogLevel)
	if err != nil {
		log.Fatal(err.Error())
	}

	hiveClient, err := thehive.GetHiveClient(shomonConf.HiveUrl, shomonConf.HiveKey, http.DefaultClient)
	if err != nil {
		log.Fatal(err.Error())
	}
	var bannerClient producer.Producer
	if shomonConf.Webhook {
		bannerClient, err = producer.GetShodanWebhook(shomonConf.ShodanKey, shomonConf.WebhookEndpoint, shomonConf.WebhookPort)
		log.Info("webhook mode is activated")
	} else {
		bannerClient, err = producer.GetShodanStreamClient(shomonConf.ShodanKey, http.DefaultClient)
		log.Info("stream listening mode is activated")
	}

	if err != nil {
		log.Fatal(err.Error())
	}

	logger.Debug(fmt.Sprintf("Config= %s", shomonConf.Print()))

	srv.Config = shomonConf
	srv.Logger = *logger
	srv.wg = wg
	srv.ctx = ctx
	srv.HiveClient = hiveClient
	srv.ShodanClient = bannerClient

	return srv, nil
}

func (s *Service) ListenStream() error {
	defer s.wg.Done()

	s.Logger.Info("starting service...")

	alertChan, err := s.ShodanClient.ListenAlerts(s.ctx)
	if err != nil {
		return err
	}

	for banner := range alertChan {
		s.wg.Add(1)
		go func() {
			id, err := s.ProcessAlert(banner)
			if err != nil {
				s.Logger.Error(err.Error())
				return
			}
			s.Logger.Info(fmt.Sprintf("Alert %s Created for %s", id, banner.IPstr))
		}()
	}

	return nil
}

func (s *Service) ProcessAlert(banner models.Service) (id string, err error) {
	defer s.wg.Done()
	s.Logger.Debug(PrintBanner(banner))

	foundService := fmt.Sprintf("%s:%d", banner.IPstr, banner.Port)
	foundServiceHash := md5.Sum([]byte(foundService))

	alert := thehive.NewAlert()
	alert.Type = s.Config.HiveType
	alert.Source = "Shodan"
	alert.SourceRef = hex.EncodeToString(foundServiceHash[:6])
	alert.Tags = s.Config.HiveTags
	alert.Title = fmt.Sprintf("Shodan Alert: %s", foundService)
	alert.AddObservable("ip", banner.IPstr)
	alert.ExternalLink = fmt.Sprintf("https://www.shodan.io/host/%s", banner.IPstr)
	alert.Description = fmt.Sprintf("[Alert Link](%s)", alert.ExternalLink)
	if s.Config.IncludeBanner {
		alert.Description = fmt.Sprintf("%s\n\n```\n\n%s\n\n```", alert.Description, PrintBanner(banner))
	}
	id, err = s.HiveClient.CreateAlert(alert)
	return id, err
}

func PrintBanner(banner models.Service) string {
	s, _ := json.MarshalIndent(banner, "", "\t")
	return string(s)
}
