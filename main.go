package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	_ "rvcmd/console"
)

//go:generate ./pckcfg.sh assets packs tools

var usetrs = true

const ua = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 Edg/123.0.0.0"

func main() {
	logrus.Infoln("RVC Models Downloader start at", time.Now().Local().Format(time.DateTime+" (MST)"))
	logrus.Infof("operating system: %s, architecture: %s", runtime.GOOS, runtime.GOARCH)
	logrus.Infoln("can use ipv6:", canUseIPv6.Get())
	ntrs := flag.Bool("notrs", false, "use standard TLS client")
	dnsf := flag.String("dns", "", "custom dns.yaml")
	cust := flag.Bool("c", false, "use custom yaml instruction")
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage:", os.Args[0], "[-notrs] [-dns dns.yaml] 'target/to/download'")
		flag.PrintDefaults()
		fmt.Println("  'target/to/download'\n        like packs/general/latest")
		fmt.Println("All available targets:")
		fmt.Println(cmdlst.String())
		return
	}
	if *ntrs {
		usetrs = false
	}
	if *dnsf != "" {
		f, err := os.Open(*dnsf)
		if err != nil {
			logrus.Errorln("open custom dns file", *dnsf, "err:", err)
			return
		}
		m := map[string][]string{}
		err = yaml.NewDecoder(f).Decode(&m)
		if err != nil {
			logrus.Errorln("decode custom dns file", *dnsf, "err:", err)
			return
		}
		_ = f.Close()
		if canUseIPv6.Get() {
			dotv6servers.add(m)
		} else {
			dotv4servers.add(m)
		}
		fmt.Println("custom dns file added")
	}
	usercfg, err := readconfig(args[0], *cust)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	err = usercfg.download(args[0], "", *cust)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	logrus.Info("all download tasks finished.")
}
