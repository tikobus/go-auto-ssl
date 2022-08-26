package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
)

var email = ""
var domain = ""
var httpPort = "20080"
var httpsPort = "20443"
var cacheDir = "/usr/local/nginx/conf/ssl"

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Useage: gossl -d domain.com [-e you@email.com] [-c ssl_dir] [-h 20080] [-s 20443]")
		return
	}

	for k, v := range os.Args {
		switch v {
		case "-e":
			email = os.Args[k+1]
		case "-d":
			domain = os.Args[k+1]
		case "-h":
			httpPort = os.Args[k+1]
		case "-s":
			httpsPort = os.Args[k+1]
		case "-c":
			cacheDir = os.Args[k+1]
		}
	}

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
		Cache:      autocert.DirCache(cacheDir),
		Email:      email,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world ssl!"))
	})

	go http.ListenAndServe(fmt.Sprintf(":%s", httpPort), certManager.HTTPHandler(nil))

	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			fmt.Println("ever 5m")
		}
	}()

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", httpsPort),
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
			NextProtos:     []string{http2.NextProtoTLS, "http/1.1"},
			MinVersion:     tls.VersionTLS12,
		},
		MaxHeaderBytes: 32 << 20,
	}

	server.ListenAndServeTLS("", "")
}
