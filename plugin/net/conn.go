package net

import (
	"fmt"
	"net"
	"time"

	"github.com/extism/go-pdk"

	"github.com/mymmrac/wasm-gate/plugin/io"
)

//go:wasmimport extism:host/user net.conn.read
func _read(connectionID int32, data uint64) int32

//go:wasmimport extism:host/user net.conn.write
func _write(connectionID int32, data uint64) int32

type Conn struct {
	connectionID int32
}

func (c *Conn) Read(b []byte) (n int, err error) {
	dataMem := pdk.Allocate(len(b))
	defer dataMem.Free()

	ioHandle := _read(c.connectionID, dataMem.Offset())
	if ioHandle < 0 {
		return 0, fmt.Errorf("failed to start read: %d", ioHandle)
	}

	readBytes := io.Ready(ioHandle)
	if readBytes < 0 {
		return 0, fmt.Errorf("failed to read: %d", readBytes)
	}

	dataMem.Load(b[:readBytes])
	return int(readBytes), nil
}

func (c *Conn) Write(b []byte) (n int, err error) {
	dataMem := pdk.AllocateBytes(b)
	defer dataMem.Free()

	ioHandle := _write(c.connectionID, dataMem.Offset())
	if ioHandle < 0 {
		return 0, fmt.Errorf("failed to start write: %d", ioHandle)
	}

	writeBytes := io.Ready(ioHandle)
	if writeBytes < 0 {
		return 0, fmt.Errorf("failed to write: %d", writeBytes)
	}

	return int(writeBytes), nil
}

func (c *Conn) Close() error {
	return nil
}

func (c *Conn) LocalAddr() net.Addr {
	return nil
}

func (c *Conn) RemoteAddr() net.Addr {
	return nil
}

func (c *Conn) SetDeadline(_ time.Time) error {
	return nil
}

func (c *Conn) SetReadDeadline(_ time.Time) error {
	return nil
}

func (c *Conn) SetWriteDeadline(_ time.Time) error {
	return nil
}
