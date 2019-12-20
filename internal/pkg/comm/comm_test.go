// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package comm

import (
	"fmt"
	"testing"
	"time"

	"github.com/gvallee/syserror/pkg/syserror"
)

const (
	namespace1 = "namespace_test1"
	namespace2 = "namespace2"
)

func sendFiniMsg(peer *PeerInfo, t *testing.T) {
	myerr := peer.SendMsg(TERMMSG, nil)
	if myerr != syserror.NoErr {
		t.Fatal("SendMsg() failed")
	}

	time.Sleep(5 * time.Second) // Give a chance for the msg to arrive before we exit which could close the connection before the test completes
}

func TestServerCreation(t *testing.T) {
	t.Log("Testing creation of a valid server ")

	// Create a server asynchronously
	server := PeerInfo{
		URL: "127.0.0.1:8888",
	}

	go server.CreateEmbeddedServer()

	// Create a simple client that will just terminate everything
	syserr := server.Connect()
	if server.conn == nil || syserr != syserror.NoErr {
		t.Fatal("cannot connect to server")
	}
	t.Log("Sending termination msg...")
	sendFiniMsg(&server, t)
}

func (info *PeerInfo) runServer(t *testing.T) {
	t.Log("Actually creating the server...")
	newPeer, mysyserr := info.CreateServer()
	if mysyserr != syserror.NoErr {
		t.Fatal("cannot create new server")
	}

	// At this point, we have a socket-level connection with a new peer
	done := 0
	t.Log("Waiting for connection handshake...")
	syserr := newPeer.HandleHandshake()
	if syserr != syserror.NoErr {
		t.Fatalf("unable to handle handshake: %s", syserr.Error())
	}
	for done != 1 {
		t.Log("Receiving data...")

		/* Wait for the termination message */
		t.Log("Waiting for the termination message...")
		msgtype, _, _, _ := newPeer.RecvMsg()
		if msgtype != TERMMSG {
			t.Fatal("received wrong type of msg")
		}

		done = 1
	}
}

func TestSendRecv(t *testing.T) {
	fmt.Print("Testing data send/recv ")
	info := PeerInfo{
		URL: "127.0.0.1:9889",
	}

	// Create a server asynchronously
	t.Log("Creating server...")
	go info.runServer(t)

	// Once we know the server is up, we connect to it
	t.Log("Server up, conencting...")
	syserr := info.Connect()
	if info.conn == nil || syserr != syserror.NoErr {
		t.Fatal("Client error: Cannot connect to server")
	}
	t.Log("Successfully connected to server")
	t.Log("Test completed, sending termination msg...")

	sendFiniMsg(&info, t)
}
