package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/elazarl/goproxy"
)

func main() {
	patches := patchTLSServerName()
	defer patches.Reset()

	listenEndpoint := flag.String("l", "127.0.0.1:8080", "proxy listen address")
	dnsEndpoint := flag.String("d", "8.8.8.8", "dns (over udp) server address")
	flag.Parse()

	dnsAddrPort, err := parseAddrPortWithDefaultPort(*dnsEndpoint, 53);
	if err != nil {
		panic(err)
	}

	net.DefaultResolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", dnsAddrPort.String())
		},
	}

	cert, err := loadCert()
	if err != nil {
		panic(err)
	}

	var customAlwaysMitm goproxy.FuncHttpsHandler = func(req string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string)  {
		return &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(cert)}, req
	};

	proxy := goproxy.NewProxyHttpServer()
	proxy.CertStore = newCertStorage()
	proxy.OnRequest().HandleConnect(customAlwaysMitm)

	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		if req.URL.Scheme == "http" {
			req.Host = "localhost"
		}
		return req, nil
	})

	log.Printf("listen at %v\n", *listenEndpoint)
	panic(http.ListenAndServe(*listenEndpoint, proxy))
}