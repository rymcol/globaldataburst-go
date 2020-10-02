package gdb

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

const (
	keepAliveTimeout = 10 * time.Second
	timeout          = 10 * time.Second
)

func newClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: keepAliveTimeout}).Dial,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}

	return client
}
