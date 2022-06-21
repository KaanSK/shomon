package conf

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

type ShomonConfig struct {
	HiveUrl          string   `koanf:"HIVE_URL"`
	HiveCaseTemplate string   `koanf:"HIVE_CASE_TEMPLATE" `
	HiveKey          string   `koanf:"HIVE_KEY"`
	HiveTags         []string `koanf:"HIVE_TAGS"`
	HiveType         string   `koanf:"HIVE_TYPE"`
	ShodanKey        string   `koanf:"SHODAN_KEY"`
	LogLevel         string   `koanf:"LOG_LEVEL"`
	IncludeBanner    bool     `koanf:"INCLUDE_BANNER"`
	Webhook          bool     `koanf:"WEBHOOK"`
	WebhookEndpoint  string   `koanf:"WEBHOOK_ENDPOINT"`
	WebhookPort      int      `koanf:"WEBHOOK_PORT"`
}

func New() (conf ShomonConfig, err error) {
	newConfig := &ShomonConfig{}
	if err := newConfig.Setup(); err != nil {
		return *newConfig, err
	}
	return *newConfig, nil
}

func (d *ShomonConfig) Setup() error {
	k := koanf.New(".")

	//Default Values
	k.Load(confmap.Provider(map[string]interface{}{
		"HIVE_URL":  "http://localhost:9000",
		"HIVE_KEY":  "NO_KEY",
		"LOG_LEVEL": "INFO",
	}, "."), nil)

	if err := k.Load(file.Provider("conf.yaml"), yaml.Parser()); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			k.Load(env.Provider("SHOMON_", ".", func(s string) string {
				return strings.TrimPrefix(s, "SHOMON_")
			}), nil)
		} else {
			return err
		}
	}

	k.Unmarshal("", &d)
	return nil
}

func (d *ShomonConfig) Print() string {
	if d != nil {
		return fmt.Sprintf("%+v", d)
	}
	return ""
}
