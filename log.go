package main

import (
	"errors"
	"io"

	"github.com/sirupsen/logrus"
)

var (
	errZeroMeterSize = errors.New("zero meter size")
)

type progressmeter struct {
	prefix string
	name   string
	size   int
	prgs   int
	lstp   int
	io.Writer
}

func newmeter(prefix, name string, size int) (pm progressmeter) {
	pm.prefix = prefix
	pm.name = name
	pm.size = size
	return
}

func (pm *progressmeter) Write(p []byte) (n int, err error) {
	if pm.size == 0 {
		return 0, errZeroMeterSize
	}
	pm.prgs += len(p)
	percent := pm.prgs * 100 / pm.size
	if percent == pm.lstp {
		return len(p), nil
	}
	logrus.Infof("%s [%2d%%] %s\t(%d/%dMB)", pm.prefix, percent, pm.name, pm.prgs/1024/1024, pm.size/1024/1024)
	pm.lstp = percent
	return len(p), nil
}
