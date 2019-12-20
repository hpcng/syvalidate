// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package blockchain

import "hash"

type Namespace struct {
	// ID is a human readable identifier for a given namespace
	ID string

	// Hash is the hash of the namespace'ID
	Hash hash.Hash
}
