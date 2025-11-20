package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	trshttp "github.com/fumiama/terasu/http"
	"github.com/pkg/errors"
)

var pcli = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	},
}

func (c *config) download(p, prefix, home, ua string, waits time.Duration, usecust, usetrs, force bool) error {
	for i, t := range c.Targets {
		if t.Refer != "" {
			refp := p[:strings.LastIndex(p, "/")+1] + t.Refer
			infof("#%s%d refer to target '%s'.", prefix, i+1, refp)
			refcfg, err := readconfig(refp, usecust)
			if err != nil {
				return err
			}
			err = refcfg.download(refp, prefix+strconv.Itoa(i+1)+".", home, ua, waits, usecust, usetrs, force)
			if err != nil {
				return err
			}
			continue
		}
		if t.OS != "" && t.OS != runtime.GOOS {
			warnf("#%s%d target required OS: %s but you are %s, skip.", prefix, i+1, t.OS, runtime.GOOS)
			continue
		}
		if t.Arch != "" && t.Arch != runtime.GOARCH {
			warnf("#%s%d target required Arch: %s but you are %s, skip.", prefix, i+1, t.Arch, runtime.GOARCH)
			continue
		}
		homefolder := path.Join(home, t.Folder)
		err := os.MkdirAll(homefolder, 0755)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("#%s%d make target folder '%s'", prefix, i+1, homefolder))
		}
		infof("#%s%d open target folder '%s'.", prefix, i+1, homefolder)
		if len(t.Copy) == 0 {
			warnf("#%s%d empty copy target.", prefix, i+1)
			continue
		}
		wg := sync.WaitGroup{}
		wg.Add(len(t.Copy))
		infof("#%s%d download copy: '%v'.", prefix, i+1, t.Copy)
		for j, cp := range t.Copy {
			go func(i int, cp, prefix string) {
				defer wg.Done()
				sleep := waits * time.Duration(i)
				if sleep > time.Millisecond {
					time.Sleep(sleep)
				}
				fname := path.Join(homefolder, cp[strings.LastIndex(cp, "/")+1:])
				if !force {
					if _, err := os.Stat(fname); err == nil || os.IsExist(err) {
						warnf("#%s%d skip exist file %s", prefix, i+1, fname)
						return
					}
				}
				req, err := http.NewRequest("GET", c.BaseURL+"/"+cp, nil)
				if err != nil {
					errorf("#%s%d new request to %s err: %v", prefix, i+1, cp, err)
					return
				}
				infof("#%s%d get: %s", prefix, i+1, req.URL)
				if len(ua) > 0 {
					req.Header.Add("user-agent", ua)
				}
				var resp *http.Response
				if usetrs {
					resp, err = trshttp.DefaultClient.Do(req)
				} else {
					resp, err = pcli.Do(req)
				}
				if err != nil {
					errorf("#%s%d get %s err: %v", prefix, i+1, req.URL, err)
					return
				}
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					err := errors.New(fmt.Sprintf("HTTP %d %s", resp.StatusCode, resp.Status))
					errorf("#%s%d get %s err: %v", prefix, i+1, req.URL, err)
					return
				}
				f, err := os.Create(fname)
				if err != nil {
					errorf("#%s%d create file %s err: %v", prefix, i+1, fname, err)
					return
				}
				defer f.Close()
				infof("#%s%d writing file %s", prefix, i+1, fname)
				pm := newmeter(fmt.Sprintf("#%s%d", prefix, i+1), fname, int(resp.ContentLength))
				defer pm.finish()
				_, err = io.Copy(io.MultiWriter(f, &pm), resp.Body)
				if err != nil {
					errorf("#%s%d download file %s err: %v", prefix, i+1, fname, err)
					return
				}
				infof("#%s%d finished download %s", prefix, i+1, fname)
			}(j, cp, fmt.Sprintf("%s%d.", prefix, i+1))
		}
		wg.Wait()
	}
	return nil
}
