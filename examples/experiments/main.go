package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

func main() {
	f2()
}

func f1() {
	conn, err := net.Dial("tcp", "example.com:80")
	if err != nil {
		panic(err)
	}

	data := make([]byte, 4096)
	n, err := conn.Read(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(n, string(data[:n]))
}

func f2() {
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		d := &net.Dialer{}
		conn, err := d.DialContext(ctx, network, addr)
		if err != nil {
			return nil, err
		}
		return &connLogger{
			conn: conn,
		}, nil
	}

	client := &http.Client{
		Transport: tr,
	}

	resp, err := client.Get("http://example.com")
	if err != nil {
		fmt.Println("Error get:", err)
		return
	}
	fmt.Println(resp)
}

type connLogger struct {
	conn net.Conn
}

func (c *connLogger) Read(b []byte) (n int, err error) {
	fmt.Println("read", len(b))
	n, err = c.conn.Read(b)
	fmt.Println("read", n, err)
	return n, err
}

func (c *connLogger) Write(b []byte) (n int, err error) {
	fmt.Println("write", len(b), "**", string(b), "**")
	n, err = c.conn.Write(b)
	fmt.Println("write", n, err)
	return n, err
}

func (c *connLogger) Close() error {
	fmt.Println("close")
	return c.conn.Close()
}

func (c *connLogger) LocalAddr() net.Addr {
	fmt.Println("local addr")
	return c.conn.LocalAddr()
}

func (c *connLogger) RemoteAddr() net.Addr {
	fmt.Println("remote addr")
	return c.conn.RemoteAddr()
}

func (c *connLogger) SetDeadline(t time.Time) error {
	fmt.Println("set deadline")
	return c.conn.SetDeadline(t)
}

func (c *connLogger) SetReadDeadline(t time.Time) error {
	fmt.Println("set read deadline")
	return c.conn.SetReadDeadline(t)
}

func (c *connLogger) SetWriteDeadline(t time.Time) error {
	fmt.Println("set write deadline")
	return c.conn.SetWriteDeadline(t)
}
