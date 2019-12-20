// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package isolated

import (
	"fmt"

	"github.com/sylabs/syvalidate/internal/pkg/blockchain"
	"github.com/sylabs/syvalidate/internal/pkg/hashcash"
)

func new() error {

	return nil
}

func chain(b blockchain.Block) error {

	return nil
}

func CreateBlock(stamp hashcash.Stamp) error {
	// Because we combine a blockchain and a file system, we need to:
	// 1. actually create a new block
	// 2. chain it to the previous block, which automatically saves the data into the FS

	block, err := blockchain.Create(stamp)
	if err != nil {
		return fmt.Errorf("failed to create block from stamp: %s", err)
	}

	// By the time the chain operation completes, the block is part of
	// the block chain. We are in isolated mode so it is still a simple
	// operation but guarantee persistency (i.e., it is in a block "cache")
	err := block.Publish()
	if err != nil {
		return fmt.Errorf("failed to chain block: %s", err)
	}

	return nil
}
