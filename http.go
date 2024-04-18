package main

// https://github.com/fumiama/comandy

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/FloatTech/ttl"
	"github.com/fumiama/terasu"
	"golang.org/x/net/http2"
)

var (
	errEmptyHostAddress = errors.New("empty host addr")
)

var httpdialer = net.Dialer{
	Timeout: time.Minute,
}

var lookupTable = ttl.NewCache[string, []string](time.Hour)

var cli = http.Client{
	Transport: &http2.Transport{
		DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
			if httpdialer.Timeout != 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, httpdialer.Timeout)
				defer cancel()
			}

			if !httpdialer.Deadline.IsZero() {
				var cancel context.CancelFunc
				ctx, cancel = context.WithDeadline(ctx, httpdialer.Deadline)
				defer cancel()
			}

			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			addrs := lookupTable.Get(host)
			if len(addrs) == 0 {
				addrs, err = resolver.LookupHost(ctx, host)
				if err != nil {
					addrs, err = net.DefaultResolver.LookupHost(ctx, host)
					if err != nil {
						return nil, err
					}
				}
				lookupTable.Set(host, addrs)
			}
			if len(addr) == 0 {
				return nil, errEmptyHostAddress
			}
			var tlsConn *tls.Conn
			for _, a := range addrs {
				conn, err := httpdialer.DialContext(ctx, network, net.JoinHostPort(a, port))
				if err != nil {
					continue
				}
				tlsConn = tls.Client(conn, cfg)
				if usetrs {
					err = terasu.Use(tlsConn).HandshakeContext(ctx)
				} else {
					err = tlsConn.HandshakeContext(ctx)
				}
				if err == nil {
					break
				}
				_ = tlsConn.Close()
			}
			return tlsConn, err
		},
	},
}
