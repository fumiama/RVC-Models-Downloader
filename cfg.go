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
	"strconv"
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

func readconfig(path string, usecust bool) (c config, err error) {
	fname := path + ".yaml"
	f, err := cfg.Open(fname)
	if usecust && err != nil {
		f, err = os.Open(fname)
	}
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

func (c *config) download(path, prefix string, usecust, force bool) error {
	for i, t := range c.Targets {
		if t.Refer != "" {
			refp := path[:strings.LastIndex(path, "/")+1] + t.Refer
			logrus.Infof("#%s%d refer to target '%s'.", prefix, i+1, refp)
			refcfg, err := readconfig(refp, usecust)
			if err != nil {
				return err
			}
			err = refcfg.download(refp, prefix+strconv.Itoa(i+1)+".", usecust, force)
			if err != nil {
				return err
			}
			continue
		}
		if t.OS != "" && t.OS != runtime.GOOS {
			logrus.Warnf("#%s%d target required OS: %s but you are %s, skip.", prefix, i+1, t.OS, runtime.GOOS)
			continue
		}
		if t.Arch != "" && t.Arch != runtime.GOARCH {
			logrus.Warnf("#%s%d target required Arch: %s but you are %s, skip.", prefix, i+1, t.Arch, runtime.GOARCH)
			continue
		}
		err := os.MkdirAll(t.Folder, 0755)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("#%s%d make target folder '%s'", prefix, i+1, t.Folder))
		}
		logrus.Infof("#%s%d open target folder '%s'.", prefix, i+1, t.Folder)
		if len(t.Copy) == 0 {
			logrus.Warningf("#%s%d empty copy target.", prefix, i+1)
			continue
		}
		wg := sync.WaitGroup{}
		wg.Add(len(t.Copy))
		logrus.Infof("#%s%d download copy: '%v'.", prefix, i+1, t.Copy)
		for j, cp := range t.Copy {
			go func(i int, cp, prefix string) {
				defer wg.Done()
				if strings.Contains(cp, "/") { // have innner folder
					infldr := t.Folder + "/" + cp[:strings.LastIndex(cp, "/")]
					err := os.MkdirAll(infldr, 0755)
					if err != nil {
						logrus.Errorf("#%s%d make target inner folder '%s' err: %v", prefix, i+1, t.Folder, err)
						return
					}
					logrus.Infof("#%s%d make target inner folder '%s'.", prefix, i+1, t.Folder)
				}
				sleep := time.Millisecond * 100 * time.Duration(i)
				if sleep > time.Millisecond {
					time.Sleep(sleep)
				}
				fname := t.Folder + "/" + cp
				if !force {
					if _, err := os.Stat(fname); err == nil || os.IsExist(err) {
						logrus.Warnf("#%s%d skip exist file %s", prefix, i+1, fname)
						return
					}
				}
				req, err := http.NewRequest("GET", c.BaseURL+"/"+cp, nil)
				if err != nil {
					logrus.Errorf("#%s%d new request to %s err: %v", prefix, i+1, cp, err)
					return
				}
				logrus.Infof("#%s%d get: %s", prefix, i+1, req.URL)
				req.Header.Add("user-agent", ua)
				resp, err := cli.Do(req)
				if err != nil {
					logrus.Errorf("#%s%d get %s err: %v", prefix, i+1, req.URL, err)
					return
				}
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					err := errors.New(fmt.Sprintf("HTTP %d %s", resp.StatusCode, resp.Status))
					logrus.Errorf("#%s%d get %s err: %v", prefix, i+1, req.URL, err)
					return
				}
				f, err := os.Create(fname)
				if err != nil {
					logrus.Errorf("#%s%d create file %s err: %v", prefix, i+1, fname, err)
					return
				}
				defer f.Close()
				logrus.Infof("#%s%d writing file %s", prefix, i+1, fname)
				pm := newmeter(fmt.Sprintf("#%s%d", prefix, i+1), fname, int(resp.ContentLength))
				_, err = io.Copy(io.MultiWriter(f, &pm), resp.Body)
				if err != nil {
					logrus.Errorf("#%s%d download file %s err: %v", prefix, i+1, fname, err)
					return
				}
				logrus.Infof("#%s%d finished download %s", prefix, i+1, fname)
			}(j, cp, fmt.Sprintf("%s%d.", prefix, i+1))
		}
		wg.Wait()
	}
	return nil
}
