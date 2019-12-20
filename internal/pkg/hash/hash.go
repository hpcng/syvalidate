// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// HashFile gets the sha256 hash for a given file
//
// todo: Code duplication from manifest, make it usable directly from manifest package
func HashFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	hasher := sha256.New()
	_, err = io.Copy(hasher, f)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(hasher.Sum(nil))
}
