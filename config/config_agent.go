package config

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"

	"github.com/opensourceways/community-robot-lib/utils"
)

type Config interface {
	Validate() error
	SetDefault()
}

// The returned value of NewConfig must be a pointer.
// Otherwise, it should deep copy the config when reading it.
type NewConfig func() Config

type ConfigAgent struct {
	mut    sync.RWMutex
	c      Config
	b      NewConfig
	md5Sum string
	t      utils.Timer
}

func NewConfigAgent(b NewConfig) ConfigAgent {
	return ConfigAgent{b: b, t: utils.NewTimer()}
}

func (ca *ConfigAgent) load(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	content := []byte(os.ExpandEnv(string(b)))
	md5Sum := fmt.Sprintf("%x", md5.Sum(content))
	if ca.md5Sum == md5Sum {
		return nil
	}

	c := ca.b()
	if err := yaml.Unmarshal(content, c); err != nil {
		return err
	}

	c.SetDefault()

	if err := c.Validate(); err != nil {
		return err
	}

	ca.mut.Lock()
	ca.c = c
	ca.md5Sum = md5Sum
	ca.mut.Unlock()

	return nil
}

func (ca *ConfigAgent) GetConfig() (string, Config) {
	ca.mut.RLock()
	c := ca.c // copy the pointer
	v := ca.md5Sum
	ca.mut.RUnlock()

	return v, c
}

// Start starts polling path for plugin config.
// If the first attempt fails, then start returns the error.
func (ca *ConfigAgent) Start(path string) error {
	if err := ca.load(path); err != nil {
		return err
	}

	l := logrus.WithField("path", path)

	ca.t.Start(
		func() {
			if err := ca.load(path); err != nil {
				l.WithError(err).Error("loading config")
			}
		},
		1*time.Minute,
		0,
	)

	return nil
}

func (ca *ConfigAgent) Stop() {
	ca.t.Stop()
}
