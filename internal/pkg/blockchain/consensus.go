// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package blockchain

import "github.com/sylabs/syvalidate/internal/pkg/comm"

// This package implements a version of the Practical Byzantine Fault Tolerance algorithm (pBFT)

type Leader struct {
	PeerInfo comm.PeerInfo
}

type Client struct {
}

func (l *Leader) HandleConsensusReq() error {
	// Get the request from the incoming message

	// Broadcasts the request to the all the nodes.

	// Generate next block

	// Send the next block back to the initiator

	return nil
}

func (c *Client) HandleNewBlockResponse() error {
	// This is one of N answers

	// Check for timeout

	// If the message is from the leader, we know how many answers we expect

	// Check the new block to see if it is consistent with what we already received

	// Save the new block data

	// If we receive all the expected answers we are all good
	// (we need R+1 answers where R is the resilience factor, i.e., the maximum number of faulty nodes allowed)

	return nil
}

func (c *Client) Consensus() error {
	// Sends a request to the leader.

	// Start a go routine to receive results

	return nil
}

func StartElection() error {
	// Prepare the request for new election

	// Send request to all know peers

	// Wait for results

	// Check results

	return nil
}
