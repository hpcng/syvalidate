// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package blockchain

import (
	"encoding/binary"
	"fmt"

	"github.com/gvallee/syserror/pkg/syserror"
	"github.com/sylabs/syvalidate/internal/pkg/cache"
	"github.com/sylabs/syvalidate/internal/pkg/comm"
)

func sendListNameSpaces(peer *comm.PeerInfo, namespaces []Namespace) error {
	// Send the number of namespace
	buff := make([]byte, 8)
	binary.LittleEndian.PutUint64(buff, uint64(len(namespaces)))
	err := peer.SendMsg(comm.DATAMSG, buff)
	if err != syserror.NoErr {
		return fmt.Errorf("failed to send the number of namespaces: %s", err.Error())
	}

	// foreach namespace, send the hash of the name
	for _, ns := range namespaces {
		buffName := ns.Hash.Sum(nil)
		err = peer.SendMsg(comm.DATAMSG, buffName)
		if err != syserror.NoErr {
			return fmt.Errorf("failed to send the namespace's hash: %s", err.Error())
		}
	}

	return nil
}

func unpackNamespace(data []byte) string {
	return string(data)
}

func recvListNamespaces(cacheBasedir string, peer *comm.PeerInfo) error {
	// The message is of the format | nb namespaces | name1 | name 2 | ... | name n |
	// A name is a sha256 hash of a string, we therefore always know its length
	msgType, size, buff, syserr := peer.RecvMsg()
	if syserr != syserror.NoErr {
		return fmt.Errorf("failed to receive the number of namespaces: %s", syserr.Error())
	}

	if msgType != comm.DATAMSG {
		return fmt.Errorf("received wrong message type")
	}

	if size != 8 {
		return fmt.Errorf("expected 8 bytes but received %d", int(size))
	}

	var namespaces []string
	numNamespaces := binary.LittleEndian.Uint64(buff)
	var i uint64
	for i = 0; i < numNamespaces; i++ {
		msgType, size, buff, syserr := peer.RecvMsg()
		if syserr != syserror.NoErr {
			return fmt.Errorf("failed to receive the number of namespaces: %s", syserr.Error())
		}

		if msgType != comm.DATAMSG {
			return fmt.Errorf("wrong message type received")
		}

		if size != 32 {
			return fmt.Errorf("expected 32 bytes but received %d", int(size))
		}

		// Unpack the namespaces'id from the message data
		ns := unpackNamespace(buff)
		namespaces = append(namespaces, ns)
	}

	// todo: add namespace to local cache
	err := cache.AddNamespaces(cacheBasedir, namespaces)
	if err != nil {
		return fmt.Errorf("failed to update cache with list of namespaces: %s", err)
	}

	return nil
}

func (l *Leader) reqNamespaceUpdate(ns string) error {
	return nil
}

func (l *Leader) handleNamespaceUpdateResp(ns string) error {
	// Post the receive

	// Mark cache as clean

	return nil
}

// NodeInit adds a new peer to the peer-to-peer network
func (l *Leader) NodeInit(cacheBasedir string) error {
	if l == nil {
		return fmt.Errorf("invalid parameter(s)")
	}

	// Connect to leader
	syserr := l.PeerInfo.Connect()
	if syserr != syserror.NoErr {
		return fmt.Errorf("failed to connect to peer %s: %s", l.PeerInfo.URL, syserr.Error())
	}

	// Load data from local cache
	namespaces, err := cache.LoadNamespaces(cacheBasedir)
	if err != nil {
		return fmt.Errorf("failed to load namespaces from cache: %s", err)
	}

	// Sync list of namespaces
	err = recvListNamespaces(cacheBasedir, &l.PeerInfo)
	if err != nil {
		return fmt.Errorf("failed to receive list of namespaces: %s", err)
	}

	// For all blockchain namespace, request the latest data from leader
	for _, ns := range namespaces {
		// todo: mark the cache for the namespace as dirty

		err = l.reqNamespaceUpdate(ns)
		if err != nil {
			return fmt.Errorf("failed to request update for namespace: %s", err)
		}

		go l.handleNamespaceUpdateResp(ns)

		// We let the update happen in the background, moving on.
		// The cache will be update as we receive the data and
		// marked as clean again.
	}

	return nil
}
