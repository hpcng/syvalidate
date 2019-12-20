// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package blockchain

import (
	"fmt"

	"github.com/sylabs/singularity-mpi/pkg/sys"
	"github.com/sylabs/syvalidate/internal/pkg/hashcash"
)

type Block struct {
	h     string
	stamp hashcash.Stamp
	prev  string
}

func (b *Block) hash() error {
	// b.hash = prev + stamp

	return nil
}

func (b *Block) setPreviousHash() error {
	// Return current previous

	// Update previous with the new hash

	return nil
}

func updatePreviousHash(hash string) error {
	return nil
}

// commitBlock makes the block immutable; this is save it to the local
// BlockFS for persistency but *not* publish it
func (b *Block) commitBlock() error {
	// Lock

	// Get previous hash
	b.setPreviousHash()

	// Hash the block to make it immutable
	err := b.hash()
	if err != nil {
		return fmt.Errorf("failed to hash stamp")
	}

	// The previous hash becomes the hash of the new block
	err = updatePreviousHash(b.h)
	if err != nil {
		return fmt.Errorf("unable to update the previous hash: %s", err)
	}

	// Persist block (which includes all the manifest from the stamp)
	err = b.LocalStore()
	if err != nil {
		return fmt.Errorf("impossible to persist block: %s", err)
	}

	// Unlock

	return nil
}

// IsolatedCreate locally creates a new block from a stamp
func IsolatedCreate(stamp hashcash.Stamp) (Block, error) {
	b := Block{
		stamp: stamp,
	}
	return b, nil
}

// ConnectedCreate initiates a consensus operation to get a new block
func ConnectedCreate(stamp hashcash.Stamp) error {
	// Send the stamp (which will actually create a block and chain it)

	return nil
}

// Publish submit the block, which will link it to the previous block
func (b *Block) Publish(sysCfg *sys.Config) error {
	// Commit the block which persists the data
	err := b.commitBlock()
	if err != nil {
		return fmt.Errorf("failed to commit block %s: %s", b.h, err)
	}

	return nil
}
