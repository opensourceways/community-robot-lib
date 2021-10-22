/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package secret implements an agent to read and reload the secrets.
package secret

import (
	"bytes"
	"sync"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/community-robot-lib/logrusutil"
	"github.com/opensourceways/community-robot-lib/utils"
)

// Agent watches a path and automatically loads the secrets stored.
type Agent struct {
	sync.RWMutex
	secretsMap map[string][]byte
	ss         []singleSecret
}

// Start creates goroutines to monitor the files that contain the secret value.
// Additionally, Start wraps the current standard logger formatter with a
// censoring formatter that removes secret occurrences from the logs.
func (a *Agent) Start(paths []string) error {
	secretsMap, err := LoadSecrets(paths)
	if err != nil {
		return err
	}

	a.secretsMap = secretsMap

	logrus.SetFormatter(logrusutil.NewCensoringFormatter(logrus.StandardLogger().Formatter, a.getSecrets))

	// Start one goroutine for each file to monitor and update the secret's values.
	for secretPath := range secretsMap {
		a.loadSingleSecret(secretPath)
	}
	return nil
}

func (a *Agent) Stop() {
	for i := range a.ss {
		a.ss[i].stop()
	}
}

// Add registers a new path to the agent.
func (a *Agent) Add(path string) error {
	secret, err := LoadSingleSecret(path)
	if err != nil {
		return err
	}

	a.setSecret(path, secret)

	// Start one goroutine for each file to monitor and update the secret's values.
	a.loadSingleSecret(path)
	return nil
}

func (a *Agent) loadSingleSecret(path string) {
	s := singleSecret{
		l:         logrus.NewEntry(logrus.StandardLogger()),
		path:      path,
		t:         utils.NewTimer(),
		setSecret: a.setSecret,
	}
	s.start()
	a.ss = append(a.ss, s)
}

// GetSecret returns the value of a secret stored in a map.
func (a *Agent) GetSecret(secretPath string) []byte {
	a.RLock()
	v := a.secretsMap[secretPath]
	a.RUnlock()

	return v
}

// setSecret sets a value in a map of secrets.
func (a *Agent) setSecret(secretPath string, secretValue []byte) {
	a.Lock()
	a.secretsMap[secretPath] = secretValue
	a.Unlock()
}

// GetTokenGenerator returns a function that gets the value of a given secret.
func (a *Agent) GetTokenGenerator(secretPath string) func() []byte {
	return func() []byte {
		return a.GetSecret(secretPath)
	}
}

const censored = "CENSORED"

var censoredBytes = []byte(censored)

// Censor replaces sensitive parts of the content with a placeholder.
func (a *Agent) Censor(content []byte) []byte {
	a.RLock()
	defer a.RUnlock()

	for _, secret := range a.secretsMap {
		content = bytes.ReplaceAll(content, secret, censoredBytes)
	}
	return content
}

func (a *Agent) getSecrets() sets.String {
	a.RLock()
	defer a.RUnlock()

	secrets := sets.NewString()
	for _, v := range a.secretsMap {
		secrets.Insert(string(v))
	}
	return secrets
}
