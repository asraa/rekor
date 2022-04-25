//
// Copyright 2021 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sharding

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/sigstore/rekor/pkg/log"
)

type LogRanges struct {
	inactive Ranges
	active   int64
}

type Ranges []LogRange

type LogRange struct {
	TreeID           int64  `yaml:"treeID"`
	TreeLength       int64  `yaml:"treeLength"`
	EncodedPublicKey string `yaml:"encodedPublicKey"`
	decodedPublicKey string
}

func NewLogRanges(path string, treeID uint) (LogRanges, error) {
	if path == "" {
		log.Logger.Info("No config file specified, skipping init of logRange map")
		return LogRanges{}, nil
	}
	if treeID == 0 {
		return LogRanges{}, errors.New("non-zero tlog_id required when passing in shard config filepath; please set the active tree ID via the `--trillian_log_server.tlog_id` flag")
	}
	// otherwise, try to read contents of the sharding config
	var ranges Ranges
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return LogRanges{}, err
	}
	if string(contents) == "" {
		log.Logger.Info("Sharding config file contents empty, skipping init of logRange map")
		return LogRanges{}, nil
	}
	if err := yaml.Unmarshal(contents, &ranges); err != nil {
		return LogRanges{}, err
	}
	for i, r := range ranges {
		decoded, err := base64.StdEncoding.DecodeString(r.EncodedPublicKey)
		if err != nil {
			return LogRanges{}, err
		}
		r.decodedPublicKey = string(decoded)
		ranges[i] = r
	}
	return LogRanges{
		inactive: ranges,
		active:   int64(treeID),
	}, nil
}

func (l *LogRanges) ResolveVirtualIndex(index int) (int64, int64) {
	indexLeft := index
	for _, l := range l.inactive {
		if indexLeft < int(l.TreeLength) {
			return l.TreeID, int64(indexLeft)
		}
		indexLeft -= int(l.TreeLength)
	}

	// If index not found in inactive trees, return the active tree
	return l.active, int64(indexLeft)
}

func (l *LogRanges) ActiveTreeID() int64 {
	return l.active
}

func (l *LogRanges) NoInactive() bool {
	return l.inactive == nil
}

// TotalInactiveLength returns the total length across all inactive shards;
// we don't know the length of the active shard.
func (l *LogRanges) TotalInactiveLength() int64 {
	var total int64
	for _, r := range l.inactive {
		total += r.TreeLength
	}
	return total
}

func (l *LogRanges) SetInactive(r []LogRange) {
	l.inactive = r
}

func (l *LogRanges) GetInactive() []LogRange {
	return l.inactive
}

func (l *LogRanges) AppendInactive(r LogRange) {
	l.inactive = append(l.inactive, r)
}

func (l *LogRanges) SetActive(i int64) {
	l.active = i
}

func (l *LogRanges) GetActive() int64 {
	return l.active
}

func (l *LogRanges) String() string {
	ranges := []string{}
	for _, r := range l.inactive {
		ranges = append(ranges, fmt.Sprintf("%d=%d", r.TreeID, r.TreeLength))
	}
	ranges = append(ranges, fmt.Sprintf("active=%d", l.active))
	return strings.Join(ranges, ",")
}

// PublicKey returns the associated public key for the given Tree ID
// and returns the active public key by default
func (l *LogRanges) PublicKey(activePublicKey, treeID string) (string, error) {
	// if no tree ID is specified, assume the active tree
	if treeID == "" {
		return activePublicKey, nil
	}
	tid, err := strconv.Atoi(treeID)
	if err != nil {
		return "", err
	}

	for _, i := range l.inactive {
		if int(i.TreeID) == tid {
			if i.decodedPublicKey != "" {
				return i.decodedPublicKey, nil
			}
			// assume the active public key if one wasn't provided
			return activePublicKey, nil
		}
	}
	if tid == int(l.active) {
		return activePublicKey, nil
	}
	return "", fmt.Errorf("%d is not a valid tree ID and doesn't have an associated public key", tid)
}