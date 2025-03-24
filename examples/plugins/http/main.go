package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/mymmrac/wasm-gate/plugin/net"
)

//go:wasmexport main
func main() {
	certPool := x509.NewCertPool()

	data, err := os.ReadFile("/etc/ssl/certs/ca-certificates.crt")
	if err != nil {
		fmt.Println("Error read cert:", err)
		return
	}

	if !certPool.AppendCertsFromPEM(data) {
		fmt.Println("Invalid cert")
		return
	}

	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.DialContext = net.DefaultDialer.DialContext
	tr.TLSClientConfig = &tls.Config{
		RootCAs: certPool,
	}

	client := &http.Client{
		Transport: tr,
	}

	start := time.Now()
	resp, err := client.Get("https://example.com")
	fmt.Println("Time:", time.Since(start))
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
