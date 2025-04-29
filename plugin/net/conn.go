package net

import (
	"fmt"
	"net"
	"time"

	"github.com/extism/go-pdk"

	"github.com/mymmrac/wape/plugin/io"
)

type Conn struct {
	connID int32
}

//go:wasmimport wape:host/env net.conn.read
func _read(connID int32, data uint64) int32

func (c *Conn) Read(b []byte) (n int, err error) {
	dataMem := pdk.Allocate(len(b))
	defer dataMem.Free()

	handle := _read(c.connID, dataMem.Offset())
	if handle < 0 {
		return 0, fmt.Errorf("failed to start read: %d", handle)
	}

	readBytes := io.Ready(handle)
	if readBytes < 0 {
		return 0, fmt.Errorf("failed to read: %d", readBytes)
	}

	dataMem.Load(b[:readBytes])
	return int(readBytes), nil
}

//go:wasmimport wape:host/env net.conn.write
func _write(connID int32, data uint64) int32

func (c *Conn) Write(b []byte) (n int, err error) {
	dataMem := pdk.AllocateBytes(b)
	defer dataMem.Free()

	handle := _write(c.connID, dataMem.Offset())
	if handle < 0 {
		return 0, fmt.Errorf("failed to start write: %d", handle)
	}

	writeBytes := io.Ready(handle)
	if writeBytes < 0 {
		return 0, fmt.Errorf("failed to write: %d", writeBytes)
	}

	return int(writeBytes), nil
}

//go:wasmimport wape:host/env net.conn.close
func _close(connID int32) int32

func (c *Conn) Close() error {
	result := _close(c.connID)
	if result < 0 {
		return fmt.Errorf("failed to close: %d", result)
	}
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
