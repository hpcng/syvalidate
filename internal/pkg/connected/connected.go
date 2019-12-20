// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package connected

import "github.com/sylabs/syvalidate/internal/pkg/hashcash"

func CreateBlock(stamp hashcash.Stamp) error {
	// Because we combine a blockchain and a file system, we need to:
	// 1. get a new block through consensus (distributed operation)
	// 2. populate the block
	// 3. save the data into the FS (distributed file system)

	return nil
}
