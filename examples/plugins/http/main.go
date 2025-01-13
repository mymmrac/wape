package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mymmrac/wasm-gate/plugin/net"
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

	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.DialContext = net.DefaultDialer.DialContext
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
