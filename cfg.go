package main

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

//go:embed cfg.zip
var cfgdata []byte

var cfg = func() *zip.Reader {
	r, err := zip.NewReader(bytes.NewReader(cfgdata), int64(len(cfgdata)))
	if err != nil {
		panic(err)
	}
	for _, f := range r.File {
		cmdlst = append(cmdlst, f.Name)
	}
	return r
}()

func readconfig(path string) (c config, err error) {
	fname := path + ".yaml"
	f, err := cfg.Open(fname)
	if err != nil {
		err = errors.Wrap(err, "invalid path")
		return
	}
	defer f.Close()
	err = yaml.NewDecoder(f).Decode(&c)
	if err != nil {
		err = errors.Wrap(err, "invalid config")
		return
	}
	return
}

type config struct {
	BaseURL string    `yaml:"BaseURL"`
	Targets []targets `yaml:"Targets"`
}

type targets struct {
	Refer  string   `yaml:"Refer"`
	Folder string   `yaml:"Folder"`
	Copy   []string `yaml:"Copy"`
	OS     string   `yaml:"OS"`
	Arch   string   `yaml:"Arch"`
}

func (c *config) download(path string) error {
	for i, t := range c.Targets {
		if t.Refer != "" {
			refp := path[:strings.LastIndex(path, "/")+1] + t.Refer
			logrus.Infof("#%d refer to target '%s'", i+1, refp)
			refcfg, err := readconfig(refp)
			if err != nil {
				return err
			}
			err = refcfg.download(refp)
			if err != nil {
				return err
			}
			continue
		}
		if t.OS != "" && t.OS != runtime.GOOS {
			logrus.Warnf("#%d target required OS: %s but you are %s, skip.", i+1, t.OS, runtime.GOOS)
			continue
		}
		if t.Arch != "" && t.Arch != runtime.GOARCH {
			logrus.Warnf("#%d target required Arch: %s but you are %s, skip.", i+1, t.Arch, runtime.GOARCH)
			continue
		}
		err := os.MkdirAll(t.Folder, 0755)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("#%d make target folder '%s'", i+1, t.Folder))
		}
		logrus.Infof("#%d open target folder '%s'", i+1, t.Folder)
		if len(t.Copy) == 0 {
			logrus.Warningf("#%d empty copy target", i+1)
			continue
		}
		wg := sync.WaitGroup{}
		wg.Add(len(t.Copy))
		logrus.Infof("#%d download copy: %v", i+1, t.Copy)
		for i, cp := range t.Copy {
			go func(i int, cp string) {
				defer wg.Done()
				sleep := time.Millisecond * 100 * time.Duration(i)
				if sleep > time.Millisecond {
					time.Sleep(sleep)
				}
				req, err := http.NewRequest("GET", c.BaseURL+"/"+cp, nil)
				if err != nil {
					logrus.Errorln("new request to", cp, "err:", err)
					return
				}
				logrus.Infoln("get:", req.URL)
				req.Header.Add("user-agent", ua)
				resp, err := cli.Do(req)
				if err != nil {
					logrus.Errorln("get", req.URL, "err:", err)
					return
				}
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					err := errors.New(fmt.Sprintf("HTTP %d %s", resp.StatusCode, resp.Status))
					logrus.Errorln("get", req.URL, "err:", err)
					return
				}
				fname := t.Folder + "/" + cp
				f, err := os.Create(fname)
				if err != nil {
					logrus.Errorln("create file", fname, "err:", err)
					return
				}
				logrus.Infoln("writing file", fname)
				defer f.Close()
				_, err = io.Copy(f, resp.Body)
				if err != nil {
					logrus.Errorln("download file", fname, "err:", err)
					return
				}
				logrus.Infoln("downloaded file", fname)
			}(i, cp)
		}
		wg.Wait()
	}
	return nil
}
