// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package hashcash

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sylabs/syvalidate/internal/pkg/hash"
)

const (
	Version = "1.01"
)

type Stamp struct {
	version   string
	bits      int
	date      time.Time // Must be of the following format: yyyy-MM-dd'T'HH:mm:ssZ
	resource  string    // IP address
	ext       string    // Key/value pair of format [name1[=val1[,val2...]];[name2[=val1[,val2...]]...]]
	rand      string
	counter   string
	signature string // The string representation of the stamp
}

// AddManifest adds an existing manifest to 'ext' since the list of manifests corresponding to the work done
// defines the SoW
func (s *Stamp) AddManifest(path string) error {
	/* A manifest is always a file that we can hash */

	// Add the hash of the file to ext
	name := filepath.Base(path)
	name = strings.TrimRight(name, ".MANIFEST")
	hash := hash.HashFile(path)
	extStr := name + "=" + hash
	if s.ext == "" {
		s.ext = extStr
	} else {
		s.ext += ";" + extStr
	}

	return nil
}

func getRandomBase64String() string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" + digits + specials
	length := 16
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})

	return string(buf)
}

func Create(ip string) Stamp {
	var s Stamp
	s.version = Version
	s.bits = 20
	s.date = time.Now()
	s.resource = ip
	s.ext = ""
	s.rand = getRandomBase64String()
	s.counter = "" // todo: properly implement counter with a binary counter, encoded in base-64 format.
	return s
}

func serializeTime(t time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02dT%02d%02d%02d-%s", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Location())
}

// Serialize creates the unique string associated to a stamp (e.g., 1:20:040806:foo::65f460d0726f420d:13a6b8)
func (s *Stamp) Serialize() string {
	return s.version + ":" + strconv.Itoa(s.bits) + ":" + serializeTime(s.date) + ":" + s.resource + ":" + s.ext + ":" + s.rand + ":" + s.counter
}

func parseTime(t string) (time.Time, error) {
	// The string we get as input looks like 2019-12-16T221815-Local which we need to transform into 2019-12-16T22:18:15-Local
	newT := t[:13]
	newT += ":"
	newT += t[13:15]
	newT += ":"
	newT += t[15:22]

	dateFormat := "yyyy-MM-ddTHHmmssZ"
	return time.Parse(dateFormat, newT)
}

func Parse(signature string) (Stamp, error) {
	var s Stamp
	var err error

	tokens := strings.Split(signature, ":")
	if len(tokens) != 7 {
		return s, HashCashWrongFormatErr
	}
	s.version = tokens[0]
	s.bits, err = strconv.Atoi(tokens[1])
	if err != nil {
		return s, HashCashWrongFormatErr
	}
	s.date, err = parseTime(tokens[2])
	if err != nil {
		log.Printf("Error while parsing date: %s", err)
		return s, HashCashWrongFormatErr
	}
	s.resource = tokens[3]
	s.ext = tokens[4]
	s.rand = tokens[5]
	s.counter = tokens[6]
	s.signature = signature

	return s, nil
}

func inSpentDatabase(s *Stamp) bool {
	return true
}

func countZeroBits(hash string) int {
	return 0
}

func (s *Stamp) IsValid() error {
	today := time.Now()

	// IF stamp.date > today + 2days THEN
	//   RETURN futuristic
	if s.date.After(today.AddDate(0, 0, 2)) {
		return HashCashFuturisticErr
	}

	// IF stamp.date < today - 28days - 2days THEN
	//   RETURN expired
	if s.date.Before(today.AddDate(0, 0, -30)) {
		return HashCashExpiredErr
	}

	// IF count_zero_bits( SHA1( stamp ) ) < 20 THEN
	//   RETURN insufficient
	h := sha1.New()
	_, err := io.WriteString(h, s.signature)
	if err != nil {
		return HashCashWrongFormatErr
	}
	hash := hex.EncodeToString(h.Sum(nil))
	if countZeroBits(hash) < 20 {
		return HashCashInsufficientErr
	}

	// IF in_spent_database( stamp ) THEN
	//   RETURN spent
	if inSpentDatabase(s) {
		return HashCashSpentErr
	}

	return nil
}

func (s *Stamp) addToSpentDatabase() error {
	return nil
}

func (s *Stamp) ValidHashCash(ip string) bool {
	/*
	   WHILE stamp = get_next_x_hashcash_header()
	     IF stamp.email == myemail THEN
	        IF valid_stamp( stamp ) THEN
	           add_to_spent_database( stamp )
	           RETURN true
	        END
	     END
	   END
	   RETURN false
	*/

	// In our case, we get a single header/stamp at a time
	if s.resource == ip {
		if s.IsValid() == nil {
			// we add the stamp to the database to avoid handling multiple times the same stamp
			err := s.addToSpentDatabase()
			if err != nil {
				return false
			}
			return true
		}
	}

	return false
}
