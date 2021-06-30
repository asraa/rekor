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

package util

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// heavily borrowed from https://github.com/google/trillian-examples/blob/master/formats/log/checkpoint.go

type Checkpoint struct {
	// Ecosystem is the ecosystem/version string
	Ecosystem string
	// Size is the number of entries in the log at this checkpoint.
	Size uint64
	// Hash is the hash which commits to the contents of the entire log.
	Hash []byte
	// OtherContent is any additional data to be included in the signed payload; each element is assumed to be one line
	OtherContent []string
}

// String returns the String representation of the Checkpoint
func (c *Checkpoint) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n%d\n%s\n", c.Ecosystem, c.Size, base64.StdEncoding.EncodeToString(c.Hash))
	for _, line := range c.OtherContent {
		fmt.Fprintf(&b, "%s\n", line)
	}
	return b.String()
}

// MarshalText returns the common format representation of this Checkpoint.
func (c *Checkpoint) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// UnmarshalText parses the common formatted checkpoint data and stores the result
// in the Checkpoint.
//
// The supplied data is expected to begin with the following 3 lines of text,
// each followed by a newline:
// <ecosystem/version string>
// <decimal representation of log size>
// <base64 representation of root hash>
// <optional non-empty line of other content>...
// <optional non-empty line of other content>...
//
// This will discard any content found after the checkpoint (including signatures)
func (c *Checkpoint) UnmarshalText(data []byte) error {
	l := bytes.Split(data, []byte("\n"))
	if len(l) < 4 {
		return errors.New("invalid checkpoint - too few newlines")
	}
	eco := string(l[0])
	if len(eco) == 0 {
		return errors.New("invalid checkpoint - empty ecosystem")
	}
	size, err := strconv.ParseUint(string(l[1]), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid checkpoint - size invalid: %w", err)
	}
	h, err := base64.StdEncoding.DecodeString(string(l[2]))
	if err != nil {
		return fmt.Errorf("invalid checkpoint - invalid hash: %w", err)
	}
	*c = Checkpoint{
		Ecosystem: eco,
		Size:      size,
		Hash:      h,
	}
	if len(l) >= 5 {
		for _, line := range l[3:] {
			if len(line) == 0 {
				break
			}
			c.OtherContent = append(c.OtherContent, string(line))
		}
	}
	return nil
}

type RekorSTH struct {
	SignedNote
}

func CheckpointValidator(strToValidate string) bool {
	c := &Checkpoint{}
	return c.UnmarshalText([]byte(strToValidate)) == nil
}

func SignedCheckpointValidator(strToValidate string) bool {
	c := SignedNote{
		Note: &Checkpoint{},
	}
	return c.UnmarshalText([]byte(strToValidate)) == nil
}

func (r *RekorSTH) SetTimestamp(timestamp uint64) {
	var ts uint64
	c := r.SignedNote.Note.(*Checkpoint)
	for i, val := range c.OtherContent {
		if n, _ := fmt.Fscanf(strings.NewReader(val), "Timestamp: %d", &ts); n == 1 {
			c.OtherContent = append(c.OtherContent[:i], c.OtherContent[i+1:]...)
		}
	}
	c.OtherContent = append(c.OtherContent, fmt.Sprintf("Timestamp: %d", timestamp))
}

func (r *RekorSTH) GetTimestamp() uint64 {
	var ts uint64
	c := r.SignedNote.Note.(*Checkpoint)
	for _, val := range c.OtherContent {
		if n, _ := fmt.Fscanf(strings.NewReader(val), "Timestamp: %d", &ts); n == 1 {
			break
		}
	}
	return ts
}

func RekorSTHValidator(strToValidate string) bool {
	r := RekorSTH{}
	return r.UnmarshalText([]byte(strToValidate)) == nil
}