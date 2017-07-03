// Copyright 2017 Dave Pederson.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Dozer: protocol module
package proto

import (
	"fmt"
	"io"
	"sync"
)

// dozer protocol specific functions.
type DozerProtocol interface {
	Init(args ...string) error
	Dial(typ, host string, port int64) (uint64, error)
	RecvFrom(id uint64, dest string, messages chan []byte, quit chan bool) error
	SendTo(id uint64, dest string, messages chan []byte, quit chan bool) error
	io.Closer
}

// protocol registry
var registry = make(map[string]DozerProtocol)

// Registry lock
var mutex = &sync.Mutex{}

// Register a protocol by name
func Register(name string, p DozerProtocol) {
	mutex.Lock()
	defer mutex.Unlock()
	if p == nil {
		panic("proto: Registered protocol is nil")
	}
	if _, dup := registry[name]; dup {
		panic("proto: Register called twice for config " + name)
	}
	registry[name] = p
}

// Load and initialize a protocol
func LoadProtocol(name string, args ...string) (DozerProtocol, error) {
	mutex.Lock()
	defer mutex.Unlock()
	p, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("proto: unknown protocol %q (forgotten import?)", name)
	}
	if err := p.Init(args...); err != nil {
		return nil, err
	}
	return p, nil
}
