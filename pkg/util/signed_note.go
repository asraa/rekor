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
	"bufio"
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/mod/sumdb/note"
)

// Interface for SignedNotes

type Note interface {
	UnmarshalText(data []byte) error
	MarshalText() ([]byte, error)
	String() string
}

type SignedNote struct {
	// A Note to sign.
	Note
	// Signatures are one or more signature lines covering the payload
	Signatures []note.Signature
}

// Sign adds a signature to a SignedCheckpoint object
// The signature is added to the signature array as well as being directly returned to the caller
func (s *SignedNote) Sign(identity string, signer crypto.Signer, opts crypto.SignerOpts) (*note.Signature, error) {
	hf := crypto.SHA256
	if opts != nil {
		hf = opts.HashFunc()
	}

	input, _ := s.Note.MarshalText()
	var digest []byte
	if hf != crypto.Hash(0) {
		hasher := hf.New()
		_, err := hasher.Write(input)
		if err != nil {
			return nil, errors.Wrap(err, "hashing checkpoint before signing")
		}
		digest = hasher.Sum(nil)
	} else {
		digest, _ = s.Note.MarshalText()
	}

	sig, err := signer.Sign(rand.Reader, digest, opts)
	if err != nil {
		return nil, errors.Wrap(err, "signing checkpoint")
	}
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(signer.Public())
	if err != nil {
		return nil, errors.Wrap(err, "marshalling public key")
	}

	pkSha := sha256.Sum256(pubKeyBytes)

	signature := note.Signature{
		Name:   identity,
		Hash:   binary.BigEndian.Uint32(pkSha[:]),
		Base64: base64.StdEncoding.EncodeToString(sig),
	}

	s.Signatures = append(s.Signatures, signature)
	return &signature, nil
}

// Verify checks that one of the signatures can be successfully verified using
// the supplied public key
func (s SignedNote) Verify(public crypto.PublicKey) bool {
	if len(s.Signatures) == 0 {
		return false
	}

	msg, _ := s.Note.MarshalText()

	//TODO: generalize this
	digest := sha256.Sum256(msg)

	for _, s := range s.Signatures {
		sigBytes, err := base64.StdEncoding.DecodeString(s.Base64)
		if err != nil {
			return false
		}
		switch pk := public.(type) {
		case *rsa.PublicKey:
			if err := rsa.VerifyPSS(pk, crypto.SHA256, digest[:], sigBytes, &rsa.PSSOptions{Hash: crypto.SHA256}); err == nil {
				return true
			}
		case *ecdsa.PublicKey:
			if ecdsa.VerifyASN1(pk, digest[:], sigBytes) {
				return true
			}
		case *ed25519.PublicKey:
			if ed25519.Verify(*pk, msg, sigBytes) {
				return true
			}
		default:
			return false
		}
	}
	return false
}

// MarshalText returns the common format representation of this SignedNote.
func (s SignedNote) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// String returns the String representation of the SignedNote
func (s SignedNote) String() string {
	var b strings.Builder
	b.WriteString(s.Note.String())
	b.WriteRune('\n')
	for _, sig := range s.Signatures {
		var hbuf [4]byte
		binary.BigEndian.PutUint32(hbuf[:], sig.Hash)
		sigBytes, _ := base64.StdEncoding.DecodeString(sig.Base64)
		b64 := base64.StdEncoding.EncodeToString(append(hbuf[:], sigBytes...))
		fmt.Fprintf(&b, "%c %s %s\n", '\u2014', sig.Name, b64)
	}

	return b.String()
}

// UnmarshalText parses the common formatted signed note data and stores the result
// in the SignedNote. THIS DOES NOT VERIFY SIGNATURES INSIDE THE CONTENT!
//
// The supplied data is expected to contain a single Note, followed by a single
// line with no comment, followed by one or more lines with the following format:
//
// \u2014 name signature
//
// * name is the string associated with the signer
// * signature is a base64 encoded string; the first 4 bytes of the decoded value is a
//   hint to the public key; it is a big-endian encoded uint32 representing the first
//   4 bytes of the SHA256 hash of the public key
func (s *SignedNote) UnmarshalText(data []byte) error {
	sc := SignedNote{
		Note: &Checkpoint{},
	}

	if err := sc.Note.UnmarshalText(data); err != nil {
		return errors.Wrap(err, "parsing checkpoint portion")
	}

	b := bufio.NewScanner(bytes.NewReader(data))
	var pastNote bool
	for b.Scan() {
		if len(b.Text()) == 0 {
			pastNote = true
			continue
		}
		if pastNote {
			var name, signature string
			if _, err := fmt.Fscanf(strings.NewReader(b.Text()), "\u2014 %s %s\n", &name, &signature); err != nil {
				return errors.Wrap(err, "parsing signature")
			}

			sigBytes, err := base64.StdEncoding.DecodeString(signature)
			if err != nil {
				return errors.Wrap(err, "decoding signature")
			}
			if len(sigBytes) < 5 {
				return errors.New("signature is too small")
			}

			sig := note.Signature{
				Name:   name,
				Hash:   binary.BigEndian.Uint32(sigBytes[0:4]),
				Base64: base64.StdEncoding.EncodeToString(sigBytes[4:]),
			}
			sc.Signatures = append(sc.Signatures, sig)
		}
	}
	if len(sc.Signatures) == 0 {
		return errors.New("no signatures found in input")
	}

	// copy sc to s
	*s = sc
	return nil
}

func SignedNoteValidator(strToValidate string) bool {
	s := SignedNote{}
	return s.UnmarshalText([]byte(strToValidate)) == nil
}