package httping

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httptrace"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/http/httpproxy"
)

type Stats struct {
	Proxy    bool
	Scheme   string
	DNS      int64
	TCP      int64
	TLS      int64
	Process  int64
	Transfer int64
	Total    int64
}

func New(address string) (Stats, error) {
	var err error
	var stats Stats
	var t0, t1, t2, t3, t4, t5, t6, t7 int64

	req, _ := http.NewRequest("GET", address, nil)
	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			t0 = time.Now().UnixNano()
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			t1 = time.Now().UnixNano()
			if info.Err != nil {
				err = info.Err
				log.Fatal(info.Err)
			}
		},
		ConnectStart: func(net, addr string) {
		},
		ConnectDone: func(net, addr string, err error) {
			if err != nil {
				log.Fatalf("unable to connect to host %v: %v", addr, err)
			}
			t2 = time.Now().UnixNano()
		},
		GotConn: func(info httptrace.GotConnInfo) {
			t3 = time.Now().UnixNano()
		},
		GotFirstResponseByte: func() {
			t4 = time.Now().UnixNano()
		},
		TLSHandshakeStart: func() {
			t5 = time.Now().UnixNano()
		},
		TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
			t6 = time.Now().UnixNano()
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	c := &http.Client{
		Timeout: 5 * time.Second,
	}
	_, err = c.Do(req)
	if err != nil {
		match, _ := regexp.MatchString("Client.Timeout exceeded", err.Error())
		if match {
			log.Fatal("Connection timeout")
		} else {
			log.Fatal(err)
		}
	}

	t7 = time.Now().UnixNano()

	if strings.HasPrefix(address, "http://") {
		stats.Scheme = "http"
	} else if strings.HasPrefix(address, "https://") {
		stats.Scheme = "https"
	}

	if t0 == 0 {
		t0 = t2
		t1 = t2
	}

	stats.DNS = t1 - t0
	stats.TCP = t2 - t1
	stats.Process = t4 - t3
	stats.Transfer = t7 - t4
	stats.TLS = t6 - t5
	stats.Total = t7 - t0

	// Detect proxies
	pc := httpproxy.FromEnvironment()
	if pc.HTTPProxy != "" {
		stats.Proxy = true
	}

	return stats, err
}
