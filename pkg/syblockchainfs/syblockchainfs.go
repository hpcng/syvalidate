// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package syblockchainfs

import (
	"fmt"

	"github.com/sylabs/singularity-mpi/pkg/sys"
	"github.com/sylabs/syvalidate/internal/pkg/connected"
	"github.com/sylabs/syvalidate/internal/pkg/hashcash"
	"github.com/sylabs/syvalidate/internal/pkg/isolated"
)

type Mode int

const (
	LeaderMode   Mode = 0
	PeerMode          = 1
	IsolatedMode Mode = 2
)

// CreateStamp is the function pointer to create a new stamp
type CreateStampFn func(string) hashcash.Stamp

// CreateBlockFn is the function pointer to create a new block
type CreateBlockFn func(hashcash.Stamp) error

type SyBlockchainFS struct {
	info Info

	sysCfg *sys.Config

	// CreateStamp is the function that creates a new stamp
	CreateStamp CreateStampFn

	// CreateBlock is the function that create a new block from a stamp
	CreateBlock CreateBlockFn
}

type Info struct {
	// Connected is set to true when the tool is expected to connect to
	// an actual blockchain. If set to false, the tool acts in a isolated
	// manner, still rely on blockchains but everything ends on the local
	// disk.
	Connected bool

	// IsLeader specifies if we are currently in the context of a leader
	// while in connected mode
	IsLeader bool
}

func initIsolatedMode(i *Info) (SyBlockchainFS, error) {
	var syBCFS SyBlockchainFS

	syBCFS.CreateStamp = hashcash.Create
	syBCFS.CreateBlock = isolated.CreateBlock

	return syBCFS, nil
}

func initConnectedMode(i *Info) (SyBlockchainFS, error) {
	var syBCFS SyBlockchainFS

	if i.IsLeader {
		err := syBCFS.Switch(LeaderMode)
		if err != nil {
			return syBCFS, fmt.Errorf("failed to switch to leader mode: %s", err)
		}
	} else {
		err := syBCFS.Switch(PeerMode)
		if err != nil {
			return syBCFS, fmt.Errorf("failed to switch to peer mode: %s", err)
		}
	}
	return syBCFS, nil
}

func (fs *SyBlockchainFS) Switch(m Mode) error {
	switch m {
	case IsolatedMode:
		return fmt.Errorf("cannot switch from isolated mode to connected mode without restart and configuration change")
	case LeaderMode:
		// A leader is not supposed to create stamps; it can only create blocks (and submit them) from stamps
		fs.CreateStamp = nil
		fs.CreateBlock = connected.CreateBlock
		fs.info.IsLeader = true
		fs.info.Connected = true
	case PeerMode:
		// A non-leader peer can only create stamps, not blocks
		fs.CreateStamp = hashcash.Create
		fs.CreateBlock = nil
		fs.info.IsLeader = false
		fs.info.Connected = true
	default:
		return fmt.Errorf("invalid mode (%d)", m)
	}

	return nil
}

func Init(i *Info) (SyBlockchainFS, error) {
	var fs SyBlockchainFS
	var err error

	if i.Connected {
		fs, err = initIsolatedMode(i)
		if err != nil {
			return fs, fmt.Errorf("failed to initialize in isolated mode: %s", err)
		}
	} else {
		fs, err = initConnectedMode(i)
		if err != nil {
			return fs, fmt.Errorf("failed to initialize in connected mode: %s", err)
		}
	}

	fs.info = *i

	return fs, nil
}
