package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jpillora/chisel/share/cio"
	"io/ioutil"
	"sync"
)

type Limit struct {
	Text string
}

type Limits struct {
	sync.RWMutex
	inner map[string]*Limits
}

type LimitsIndex struct {
	*cio.Logger
	*Limits
	configFile string
}

func (u *LimitsIndex) loadLimitIndex() error {
	if u.configFile == "" {
		return errors.New("configuration file not set")
	}
	b, err := ioutil.ReadFile(u.configFile)
	if err != nil {
		return fmt.Errorf("Failed to read limit file: %s, error: %s", u.configFile, err)
	}
	var raw map[string][]string
	if err := json.Unmarshal(b, &raw); err != nil {
		return errors.New("Invalid JSON: " + err.Error())
	}
	var limits []*Limit
	for l := range raw {
		limit := &Limit{}
		limit.Text = ParseLimit(l)
		limits = append(limits, limit)
	}

	return nil
}

// ParseLimit TODO
func ParseLimit(l string) string {
	return ""
}

func (u *LimitsIndex) LoadLimits(configFile string) error {
	u.configFile = configFile
	u.Infof("Loading configuration file %s", configFile)
	if err := u.loadLimitIndex(); err != nil {
		return err
	}
	return nil
}
