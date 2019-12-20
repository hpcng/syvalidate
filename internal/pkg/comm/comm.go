// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package comm

import (
	"encoding/binary"
	"log"
	"net"
	"time"

	"github.com/gvallee/syserror/pkg/syserror"
)

// Messages types
const (
	// INVALID represents an invalid msg
	INVALID = "INVA"
	// TERMMSG is a termination message
	TERMMSG = "TERM"
	// CONNREQ is a connection request (initiate a connection handshake
	CONNREQ = "CONN"
	// CONNACK is a response to a connection request (Connection ack)
	CONNACK = "CACK"
	// DATAMSG represents a data msg
	DATAMSG = "DATA"
)

// Structure to store server information (host we connect to)
type PeerInfo struct {
	conn    net.Conn
	timeout int

	// URL is the IP/port to use to connect to the peer
	URL string
}

// Receive and parse a message header (4 character)
func (p *PeerInfo) GetHeader() (string, syserror.SysError) {
	if p.conn == nil {
		return INVALID, syserror.ErrNotAvailable
	}

	hdr := make([]byte, 4)
	if hdr == nil {
		return INVALID, syserror.ErrOutOfRes
	}

	/* read the msg type */
	s, err := p.conn.Read(hdr)
	if err != nil {
		log.Println("ERROR:", err.Error())
		return INVALID, syserror.ErrFatal
	}
	// Connection is closed
	if s == 0 {
		log.Println("Connection closed")
		return TERMMSG, syserror.NoErr
	}
	// Read returns an error (we test this only as a second test since s=0 means socket closed,
	// which returns also EOF as an error but we want to handle it as normal termination, not an error case
	if err != nil {
		log.Println("ERROR:", err.Error())
		return INVALID, syserror.ErrFatal
	}

	if s > 0 && s != 4 {
		log.Println("ERROR Cannot recv header")
		return INVALID, syserror.NoErr
	}

	// Disconnect request
	if string(hdr[:s]) == TERMMSG {
		log.Println("Recv'd disconnect request")
		return TERMMSG, syserror.NoErr
	}

	if string(hdr[:s]) == CONNREQ {
		log.Println("Recv'd connection request")
		return CONNREQ, syserror.NoErr
	}

	if string(hdr[:s]) == CONNACK {
		log.Println("Recv'd connection ACK")
		return CONNACK, syserror.NoErr
	}

	log.Println("Invalid msg header")
	return INVALID, syserror.ErrFatal
}

func (p *PeerInfo) getPayloadSize() (uint64, syserror.SysError) {
	if p == nil || p.conn == nil {
		return 0, syserror.ErrNotAvailable
	}

	ps := make([]byte, 8) // Payload size is always 8 bytes
	s, myerr := p.conn.Read(ps)
	if s != 8 || myerr != nil {
		log.Println("ERROR: expecting 8 bytes but received", s)
		return 0, syserror.ErrFatal
	}

	return binary.LittleEndian.Uint64(ps), syserror.NoErr
}

func (p *PeerInfo) getPayload(size uint64) ([]byte, syserror.SysError) {
	if p == nil || p.conn == nil {
		return nil, syserror.ErrNotAvailable
	}

	payload := make([]byte, size)
	s, myerr := p.conn.Read(payload)
	if uint64(s) != size || myerr != nil {
		log.Println("ERROR: expecting ", size, "but received", s)
		return nil, syserror.ErrFatal
	}

	return payload, syserror.NoErr
}

func (p *PeerInfo) handleConnReq(size uint64, payload []byte) syserror.SysError {
	if p == nil {
		log.Println("ERROR: local server is not initialized")
		return syserror.ErrFatal
	}

	// Send CONNACK with the payload
	log.Println("Sending CONNACK...")
	err := p.SendMsg(CONNACK, payload)
	if err != syserror.NoErr {
		return err
	}

	return syserror.NoErr
}

// HandleHandshake receives and handles a CONNREQ message, i.e., a client trying to connect
func (p *PeerInfo) HandleHandshake() syserror.SysError {
	/* Handle the CONNREQ message */
	msgtype, payload_size, payload, err := p.RecvMsg()
	if err != syserror.NoErr {
		return syserror.ErrFatal
	}

	if msgtype == CONNREQ {
		err := p.handleConnReq(payload_size, payload)
		if err != syserror.NoErr {
			return err
		}
	} else {
		return syserror.NoErr // We did not get the expected CONNREQ msg
	}

	return syserror.NoErr
}

func (p *PeerInfo) sendMsgType(msgType string) syserror.SysError {
	if p.conn == nil {
		return syserror.ErrFatal
	}

	s, err := p.conn.Write([]byte(msgType))
	if s == 0 || err != nil {
		return syserror.ErrFatal
	}

	return syserror.NoErr
}

// SendMsg sends a basic message
func (p *PeerInfo) SendMsg(msgType string, payload []byte) syserror.SysError {
	if p.conn == nil {
		return syserror.ErrFatal
	}

	hdrerr := p.sendMsgType(msgType)
	if hdrerr != syserror.NoErr {
		return hdrerr
	}

	if payload != nil {
		syserr := p.sendUint64(uint64(len(payload)))
		if syserr != syserror.NoErr {
			return syserr
		}
		s, err := p.conn.Write(payload)
		if s != len(payload) || err != nil {
			log.Println("[ERROR] write operation failed")
			return syserror.ErrFatal
		}
	} else {
		syserr := p.sendUint64(uint64(0))
		if syserr != syserror.NoErr {
			return syserr
		}
	}

	return syserror.NoErr
}

func (p *PeerInfo) sendUint64(value uint64) syserror.SysError {
	if p == nil || p.conn == nil {
		return syserror.ErrFatal
	}

	buff := make([]byte, 8)
	if buff == nil {
		return syserror.ErrOutOfRes
	}

	binary.LittleEndian.PutUint64(buff, value)
	s, myerr := p.conn.Write(buff)
	if myerr != nil {
		log.Println(myerr.Error())
	}
	if myerr == nil && s != 8 {
		log.Println("Received ", s, "bytes, instead of 8")
	}
	if s != 8 || myerr != nil {
		return syserror.ErrFatal
	}

	return syserror.NoErr
}

// RecvMsg receives a basic message
func (p *PeerInfo) RecvMsg() (string, uint64, []byte, syserror.SysError) {
	if p == nil || p.conn == nil {
		return "", 0, nil, syserror.ErrFatal
	}

	msgtype, err := p.GetHeader()
	// Messages without payload
	if msgtype == "TERM" || msgtype == "INVA" || err != syserror.NoErr {
		return msgtype, 0, nil, syserror.ErrFatal
	}

	// Get the payload size
	payload_size, err := p.getPayloadSize()
	if err != syserror.NoErr {
		return msgtype, 0, nil, syserror.ErrFatal
	}
	if payload_size == 0 {
		return msgtype, 0, nil, syserror.NoErr
	}

	// Get the payload
	buff, err := p.getPayload(payload_size)
	if err != syserror.NoErr {
		return msgtype, payload_size, buff, syserror.NoErr
	}

	return msgtype, payload_size, buff, syserror.NoErr
}

// ConnectHandshake initiates a connection handshake
func (p *PeerInfo) ConnectHandshake() syserror.SysError {
	err := p.SendMsg(CONNREQ, nil)
	if err != syserror.NoErr {
		return err
	}

	// Receive the CONNACK, the payload is the block_sizw
	msgtype, size, _, recverr := p.RecvMsg()
	if recverr != syserror.NoErr || msgtype != CONNACK || size < 0 {
		return syserror.ErrFatal
	}
	// todo: handle payload

	return syserror.NoErr
}

func (p *PeerInfo) connectHandshake() syserror.SysError {
	err := p.SendMsg(CONNREQ, nil)
	if err != syserror.NoErr {
		return err
	}

	// Receive the CONNACK, the payload is the block_sizw
	msgtype, s, _, err := p.RecvMsg()
	if err != syserror.NoErr || msgtype != CONNACK || s < 0 {
		return syserror.ErrFatal
	}

	return syserror.NoErr
}

func (p *PeerInfo) Connect() syserror.SysError {
	retry := 0
	if p == nil {
		return syserror.ErrFatal
	}

	var err error
Retry:
	p.conn, err = net.Dial("tcp", p.URL)
	if err != nil {
		log.Printf("Dial failed: %s", err.Error())
		if retry < 5 {
			retry++
			log.Printf("Retrying after %d seconds\n", retry)
			time.Sleep(time.Duration(retry) * time.Second)
			goto Retry
		}
		return syserror.ErrOutOfRes
	}

	syserr := p.connectHandshake()
	if syserr != syserror.NoErr {
		return syserror.ErrFatal
	}

	return syserror.NoErr
}

func (info *PeerInfo) doServer() syserror.SysError {
	done := 0

	syserr := info.HandleHandshake()
	if syserr != syserror.NoErr {
		log.Printf("[ERROR] handshake with client failed: %s", syserr.Error())
	}
	for done != 1 {
		// Handle the message header
		msgtype, syserr := info.GetHeader()
		if syserr != syserror.NoErr {
			done = 1
		}
		if msgtype == "INVA" || msgtype == "TERM" {
			done = 1
		}

		//  Handle the paylaod size
		payload_size, pserr := info.getPayloadSize()
		if pserr != syserror.NoErr {
			done = 1
		}

		// Handle the payload when necessary
		var payload []byte = nil
		if payload_size != 0 {
			payload, syserr = info.getPayload(payload_size)
		}
		if payload_size != 0 && payload != nil {
			return syserror.ErrFatal
		}
	}

	info.conn.Close()
	log.Printf("Go routine for peer %s done\n", info.URL)

	return syserror.NoErr
}

func (info *PeerInfo) CreateEmbeddedServer() syserror.SysError {
	if info == nil {
		return syserror.ErrFatal
	}

	log.Println("Creating embedded server...")
	listener, err := net.Listen("tcp", info.URL)
	if err != nil {
		return syserror.ErrFatal
	}

	log.Println("Server created on", info.URL)

	var conn net.Conn
	for {
		conn, err = listener.Accept()
		log.Println("Connection accepted")
		var newPeer PeerInfo
		newPeer.conn = conn
		if err != nil {
			log.Println("[ERROR] ", err.Error())
		}

		log.Println("Creating new Go routine for new peer...")
		go newPeer.doServer()
	}

	// todo: fix termination
	//return syserror.NoErr
}

func (info *PeerInfo) CreateServer() (PeerInfo, syserror.SysError) {
	var newPeer PeerInfo
	if info == nil {
		return newPeer, syserror.ErrFatal
	}

	listener, err := net.Listen("tcp", info.URL)
	if err != nil {
		log.Printf("failed to listen on socket: %s", err)
		return newPeer, syserror.ErrFatal
	}

	if info.timeout > 0 {
		listener.(*net.TCPListener).SetDeadline(time.Now().Add(time.Duration(info.timeout) * time.Second))
	}

	newPeer.conn, err = listener.Accept()
	if err != nil {
		return newPeer, syserror.ErrFatal
	}

	return newPeer, syserror.NoErr
}
