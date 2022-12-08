package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jpillora/chisel/share/cio"
	"io/ioutil"
	"reflect"
	"sync"
)

type Limits struct {
	sync.RWMutex
	inner []*Remote
}

type LimitsIndex struct {
	*cio.Logger
	*Limits
	configFile string
}

func NewLimits() *Limits {
	return &Limits{inner: []*Remote{}}
}

func NewLimitsIndex(logger *cio.Logger) *LimitsIndex {
	return &LimitsIndex{
		Logger: logger.Fork("limits"),
		Limits: NewLimits(),
	}
}

// Len returns the numbers of limits
func (u *Limits) Len() int {
	u.RLock()
	l := len(u.inner)
	u.RUnlock()
	return l
}

func (l *Limits) In(limit Remote) bool {
	l.RLock()
	same := false
	for _, remote := range l.inner {
		same = reflect.DeepEqual(remote, &limit)
	}
	l.RUnlock()
	return same

}

func (u *Limits) Reset(limits []*Remote) {
	m := []*Remote{}
	for _, u := range limits {
		m = append(m, u)
	}
	u.Lock()
	u.inner = m
	u.Unlock()
}

func (u *LimitsIndex) loadLimitIndex() error {
	if u.configFile == "" {
		return errors.New("configuration file not set")
	}
	b, err := ioutil.ReadFile(u.configFile)
	if err != nil {
		return fmt.Errorf("failed to read limit file: %s, error: %s", u.configFile, err)
	}
	var raw []string
	if err := json.Unmarshal(b, &raw); err != nil {
		return errors.New("Invalid JSON: " + err.Error())
	}
	var remotes []*Remote
	for _, l := range raw {
		remote, err := DecodeRemote(l)
		u.Infof("Limit: %s", l)

		if err != nil {
			return fmt.Errorf("failed to decode the remote string: %s", err)
		}
		remotes = append(remotes, remote)
	}
	u.Reset(remotes)
	return nil
}

func (u *LimitsIndex) LoadLimits(configFile string) error {
	u.configFile = configFile
	u.Infof("Loading configuration file %s", configFile)
	if err := u.loadLimitIndex(); err != nil {
		return err
	}
	return nil
}
