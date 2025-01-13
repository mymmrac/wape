package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/extism/go-pdk"
)

//go:wasmexport main
func main() {
	certPool := x509.NewCertPool()

	data, err := os.ReadFile("/etc/ssl/certs/ca-certificates.crt")
	if err != nil {
		panic(err)
	}

	if !certPool.AppendCertsFromPEM(data) {
		panic("append certificates")
	}

	d := &dialer{}

	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.DialContext = d.DialContext
	tr.TLSClientConfig = &tls.Config{
		RootCAs: certPool,
	}

	client := &http.Client{
		Transport: tr,
	}

	resp, err := client.Get("https://example.com")
	if err != nil {
		fmt.Println("Error get:", err)
		return
	}
	fmt.Println("Status:", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error read:", err)
		return
	}

	fmt.Println("Body:", string(body))
}

type dialer struct{}

//go:wasmimport extism:host/user dial
func _dial(network, addr uint64) int32

func (d *dialer) DialContext(_ context.Context, network, addr string) (net.Conn, error) {
	networkMem := pdk.AllocateString(network)
	defer networkMem.Free()

	addrMem := pdk.AllocateString(addr)
	defer addrMem.Free()

	connectionID := _dial(networkMem.Offset(), addrMem.Offset())
	if connectionID < 0 {
		return nil, fmt.Errorf("failed to dial: %d", connectionID)
	}

	return &conn{
		connectionID: connectionID,
	}, nil
}

type conn struct {
	connectionID int32
}

//go:wasmimport extism:host/user read
func _read(connectionID int32, data uint64) int32

func (c *conn) Read(b []byte) (n int, err error) {
	dataMem := pdk.Allocate(len(b))
	defer dataMem.Free()

	ioHandle := _read(c.connectionID, dataMem.Offset())
	if ioHandle < 0 {
		return 0, fmt.Errorf("failed to start read: %d", ioHandle)
	}

	readBytes := ioReady(ioHandle)
	if readBytes < 0 {
		return 0, fmt.Errorf("failed to read: %d", readBytes)
	}

	dataMem.Load(b[:readBytes])
	return int(readBytes), nil
}

//go:wasmimport extism:host/user write
func _write(connectionID int32, data uint64) int32

func (c *conn) Write(b []byte) (n int, err error) {
	dataMem := pdk.AllocateBytes(b)
	defer dataMem.Free()

	ioHandle := _write(c.connectionID, dataMem.Offset())
	if ioHandle < 0 {
		return 0, fmt.Errorf("failed to start write: %d", ioHandle)
	}

	writeBytes := ioReady(ioHandle)
	if writeBytes < 0 {
		return 0, fmt.Errorf("failed to write: %d", writeBytes)
	}

	return int(writeBytes), nil
}

//go:wasmimport extism:host/user io_ready
func _ioReady(ioHandle int32) int32

func ioReady(ioHandle int32) int32 {
	var result int32
	for {
		result = _ioReady(ioHandle)
		if result != 0 {
			return result
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func (c *conn) Close() error {
	fmt.Println("Close")
	return nil
}

func (c *conn) LocalAddr() net.Addr {
	return nil
}

func (c *conn) RemoteAddr() net.Addr {
	return nil
}

func (c *conn) SetDeadline(_ time.Time) error {
	return nil
}

func (c *conn) SetReadDeadline(_ time.Time) error {
	return nil
}

func (c *conn) SetWriteDeadline(_ time.Time) error {
	return nil
}
