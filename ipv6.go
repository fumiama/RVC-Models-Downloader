package main

// https://github.com/fumiama/comandy

import (
	"context"
	"net/http"
	"time"

	"github.com/RomiChan/syncx"
)

var canUseIPv6 = syncx.Lazy[bool]{Init: func() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://v6.ipv6-test.com/json/widgetdata.php?callback=?", nil)
	if err != nil {
		return false
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	_ = resp.Body.Close()
	return true
}}
