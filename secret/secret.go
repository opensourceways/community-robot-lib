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

package secret

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/opensourceways/community-robot-lib/utils"
	"github.com/sirupsen/logrus"
)

// LoadSecrets loads multiple paths of secrets and add them in a map.
func LoadSecrets(paths []string) (map[string][]byte, error) {
	secretsMap := make(map[string][]byte, len(paths))

	for _, path := range paths {
		secretValue, err := LoadSingleSecret(path)
		if err != nil {
			return nil, err
		}
		secretsMap[path] = secretValue
	}
	return secretsMap, nil
}

// LoadSingleSecret reads and returns the value of a single file.
func LoadSingleSecret(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", path, err)
	}
	return bytes.TrimSpace(b), nil
}

type singleSecret struct {
	l    *logrus.Entry
	path string

	t           utils.Timer
	skips       int
	lastModTime time.Time
	setSecret   func(string, []byte)
}

// reloadSecret will begin polling the secret file at the path. If the first load
// fails, Start with return the error and abort. Future load failures will log
// the failure message but continue attempting to load.
func (s *singleSecret) load() {
	if s.skips < 600 {
		// Check if the file changed to see if it needs to be re-read.
		secretStat, err := os.Stat(s.path)
		if err != nil {
			//l.WithField("secret-path", s.path).WithError(err).Error("Error loading secret file.")
			s.l.WithError(err).Error("Error statting secret file.")
			return
		}

		t := secretStat.ModTime()
		if !t.After(s.lastModTime) {
			s.skips++
			return // file hasn't been modified
		}

		s.lastModTime = t
	}

	secretValue, err := LoadSingleSecret(s.path)
	if err != nil {
		s.l.WithError(err).Error("Error loading secret.")
		return
	}

	s.setSecret(s.path, secretValue)
	s.skips = 0
}

func (s singleSecret) start() {
	s.t.Start(s.load, 1*time.Second)
}

func (s singleSecret) stop() {
	s.t.Stop()
}
