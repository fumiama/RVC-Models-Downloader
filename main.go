package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/fumiama/terasu/dns"
	"github.com/fumiama/terasu/ip"
	ui "github.com/gizak/termui/v3"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	_ "rvcmd/console"
)

//go:generate ./pckcfg.sh assets

var (
	notui = false
	sc    screen
)

func main() {
	ntrs := flag.Bool("notrs", false, "use standard TLS client")
	dnsf := flag.String("dns", "", "custom dns.yaml")
	cust := flag.Bool("c", false, "use custom yaml instruction")
	force := flag.Bool("f", false, "force download even file exists")
	wait := flag.Uint("w", 4, "connection waiting seconds")
	ua := flag.String("ua", defua, "customize user agent")
	h := flag.Bool("h", false, "display this help")
	home := flag.String("H", ".", "download under this path")
	flag.BoolVar(&notui, "notui", false, "use plain text instead of TUI")
	flag.BoolVar(&ip.IsIPv6Available, "6", false, "use ipv6")
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 || *h {
		fmt.Println("Usage:", os.Args[0], "[-notrs] [-dns dns.yaml] 'target/to/download'")
		flag.PrintDefaults()
		fmt.Println("  'target/to/download'\n        like packs/general/latest")
		fmt.Println("All available targets:")
		fmt.Println(cmdlst.String())
		return
	}
	err := os.MkdirAll(*home, 0755)
	if err != nil {
		logrus.Errorln("mkdirs of path", *home, "err:", err)
		return
	}
	if notui {
		logrus.Infoln("RVC Models Downloader start at", time.Now().Local().Format(time.DateTime+" (MST)"))
		logrus.Infof("operating system: %s, architecture: %s", runtime.GOOS, runtime.GOARCH)
		logrus.Infoln("is ipv6 available:", ip.IsIPv6Available)
	} else {
		if err := ui.Init(); err != nil {
			logrus.Errorln("failed to initialize termui:", err)
			return
		}
		defer ui.Close()
		sc = newscreen()
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if *dnsf != "" {
			f, err := os.Open(*dnsf)
			if err != nil {
				errorln("open custom dns file", *dnsf, "err:", err)
				return
			}
			m := dns.Config{}
			err = yaml.NewDecoder(f).Decode(&m)
			if err != nil {
				errorln("decode custom dns file", *dnsf, "err:", err)
				return
			}
			_ = f.Close()
			if ip.IsIPv6Available {
				dns.IPv6Servers.Add(&m)
			} else {
				dns.IPv4Servers.Add(&m)
			}
			infoln("custom dns file added")
		}
		usercfg, err := readconfig(args[0], *cust)
		if err != nil {
			errorln(err)
			return
		}
		ch := make(chan struct{})
		go func() {
			err := usercfg.download(args[0], "", *home, *ua, time.Second*time.Duration(*wait), *cust, !*ntrs, *force)
			ch <- struct{}{}
			if err != nil {
				errorln(err)
				return
			}
		}()
		select {
		case <-ch:
			infoln("all download tasks finished.")
		case <-ctx.Done():
			logrus.Warnln("download canceled")
		}
	}()
	if notui {
		wg.Wait()
	} else {
		sc.show(time.Second)
	}
}
