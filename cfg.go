package main

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

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
