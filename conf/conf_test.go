package conf

import (
	"testing"
)

func TestGetConf(t *testing.T) {
	conf, err := New()
	if err != nil {
		t.Errorf(err.Error())
	}
	if conf.LogLevel == "" {
		t.Errorf("Config could not be populated")
	}
}
