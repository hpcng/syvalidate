// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package blockchain

// LocalStore stores the blockchain locally using our
// BlockFS file system
func (b *Block) LocalStore() error {
	// Write all manifests to the FS
	sz, err = fs.Write()

	// Write the block itself
	sz, err = fs.Write()

	return nil
}
