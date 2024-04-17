package main

// https://github.com/fumiama/comandy

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/fumiama/terasu"
)

var (
	errNoDNSAvailable = errors.New("no dns available")
)

var dnsdialer = net.Dialer{
	Timeout: time.Second * 8,
}

type dnsstat struct {
	A string
	E bool
}

type dnsservers struct {
	sync.RWMutex
	m map[string][]*dnsstat
}

// hasrecord no lock, use under lock
func hasrecord(lst []*dnsstat, a string) bool {
	for _, addr := range lst {
		if addr.A == a {
			return true
		}
	}
	return false
}

func (ds *dnsservers) add(m map[string][]string) {
	ds.Lock()
	defer ds.Unlock()
	addList := map[string][]*dnsstat{}
	for host, addrs := range m {
		for _, addr := range addrs {
			if !hasrecord(ds.m[host], addr) && !hasrecord(addList[host], addr) {
				addList[host] = append(addList[host], &dnsstat{addr, true})
			}
		}
	}
	for host, addrs := range addList {
		ds.m[host] = append(ds.m[host], addrs...)
	}
}

func (ds *dnsservers) dial(ctx context.Context) (tlsConn *tls.Conn, err error) {
	err = errNoDNSAvailable

	ds.RLock()
	defer ds.RUnlock()

	if dnsdialer.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, dnsdialer.Timeout)
		defer cancel()
	}

	if !dnsdialer.Deadline.IsZero() {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, dnsdialer.Deadline)
		defer cancel()
	}

	var conn net.Conn
	for host, addrs := range ds.m {
		for _, addr := range addrs {
			if !addr.E {
				continue
			}
			conn, err = dnsdialer.DialContext(ctx, "tcp", addr.A)
			if err != nil {
				addr.E = false // no need to acquire write lock
				continue
			}
			tlsConn = tls.Client(conn, &tls.Config{ServerName: host})
			if usetrs {
				err = terasu.Use(tlsConn).HandshakeContext(ctx)
			} else {
				err = tlsConn.HandshakeContext(ctx)
			}
			if err == nil {
				return
			}
			_ = tlsConn.Close()
			addr.E = false // no need to acquire write lock
		}
	}
	return
}

var dotv6servers = dnsservers{
	m: map[string][]*dnsstat{
		"dot.sb": {
			{"[2a09::]:853", true},
			{"[2a11::]:853", true},
		},
		"dns.google": {
			{"[2001:4860:4860::8888]:853", true},
			{"[2001:4860:4860::8844]:853", true},
		},
		"cloudflare-dns.com": {
			{"[2606:4700:4700::1111]:853", true},
			{"[2606:4700:4700::1001]:853", true},
		},
		"dns.umbrella.com": {
			{"[2620:0:ccc::2]:853", true},
			{"[2620:0:ccd::2]:853", true},
		},
		"dns10.quad9.net": {
			{"[2620:fe::10]:853", true},
			{"[2620:fe::fe:10]:853", true},
		},
	},
}

var dotv4servers = dnsservers{
	m: map[string][]*dnsstat{
		"dot.sb": {
			{"185.222.222.222:853", true},
			{"45.11.45.11:853", true},
		},
		"dns.google": {
			{"8.8.8.8:853", true},
			{"8.8.4.4:853", true},
		},
		"cloudflare-dns.com": {
			{"1.1.1.1:853", true},
			{"1.0.0.1:853", true},
		},
		"dns.umbrella.com": {
			{"208.67.222.222:853", true},
			{"208.67.220.220:853", true},
		},
		"dns10.quad9.net": {
			{"9.9.9.10:853", true},
			{"149.112.112.10:853", true},
		},
	},
}

var resolver = &net.Resolver{
	PreferGo: true,
	Dial: func(ctx context.Context, _, _ string) (net.Conn, error) {
		if canUseIPv6.Get() {
			return dotv6servers.dial(ctx)
		}
		return dotv4servers.dial(ctx)
	},
}
