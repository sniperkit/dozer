// Copyright 2015 Dave Pederson.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Dozer: main module
package dozer

import (
	"errors"
	"github.com/zdavep/dozer/proto"
	_ "github.com/zdavep/dozer/proto/stomp"
	_ "github.com/zdavep/dozer/proto/zmq4"
)

// Supported messaging protocols.
var validProto = map[string]bool{
	"stomp": true,
	"zmq4":  true,
}

// Core dozer type.
type Dozer struct {
	dest      string
	protoName string
	ctxType   string
	protocol  proto.DozerProtocol
}

// Create a new Dozer queue.
func Queue(queue string) *Dozer {
	return &Dozer{dest: queue}
}

// Create a new Dozer topic.
func Topic(topic string) *Dozer {
	return &Dozer{dest: "/topic/" + topic}
}

// Create a new Dozer socket.
func Socket(typ string) *Dozer {
	return &Dozer{ctxType: typ}
}

// Set the message type field
func (d *Dozer) WithMessageType(typ string) *Dozer {
	d.ctxType = typ
	return d
}

// Set the protocol name field
func (d *Dozer) WithProtocol(protocolName string) *Dozer {
	d.protoName = protocolName
	return d
}

// Syntactic sugar - calls Connect.
func (d *Dozer) Bind(host string, port int64) error {
	return d.Connect(host, port)
}

// Connect or bind to a host and port.
func (d *Dozer) Connect(host string, port int64) error {
	if _, ok := validProto[d.protoName]; !ok {
		return errors.New("Unsupported protocol")
	}
	p, err := proto.LoadProtocol(d.protoName, d.ctxType)
	if err != nil {
		return err
	}
	if err := p.Dial(host, port); err != nil {
		return err
	}
	d.protocol = p
	return nil
}

// Receive messages from the lower level protocol and forward them to a channel until a quit signal fires.
func (d *Dozer) RecvLoop(messages chan []byte, quit chan bool) error {
	if err := d.protocol.RecvFrom(d.dest); err != nil {
		return err
	}
	if err := d.protocol.RecvLoop(messages, quit); err != nil {
		return err
	}
	return nil
}

// Send messages to the lower level protocol from a channel until a quit signal fires.
func (d *Dozer) SendLoop(messages chan []byte, quit chan bool) error {
	if err := d.protocol.SendTo(d.dest); err != nil {
		return err
	}
	if err := d.protocol.SendLoop(messages, quit); err != nil {
		return err
	}
	return nil
}
