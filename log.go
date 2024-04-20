package main

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gizak/termui/v3/widgets"
	"github.com/sirupsen/logrus"
)

var (
	errZeroMeterSize = errors.New("zero meter size")
)

func infof(format string, args ...any) {
	if notui {
		logrus.Infof(format, args...)
	} else {
		sc.infof(format, args...)
	}
}

func (s *screen) infof(format string, args ...any) {
	s.Lock()
	defer s.Unlock()
	s.logroll.Rows = append(s.logroll.Rows, fmt.Sprintf("[INFO](fg:blue) "+format, args...))
	s.logroll.ScrollDown()
}

func warnf(format string, args ...any) {
	if notui {
		logrus.Warnf(format, args...)
	} else {
		sc.warnf(format, args...)
	}
}

func (s *screen) warnf(format string, args ...any) {
	s.Lock()
	defer s.Unlock()
	s.logroll.Rows = append(s.logroll.Rows, fmt.Sprintf("[WARN](fg:yellow) "+format, args...))
	s.logroll.ScrollDown()
}

func errorf(format string, args ...any) {
	if notui {
		logrus.Errorf(format, args...)
	} else {
		sc.errorf(format, args...)
	}
}

func (s *screen) errorf(format string, args ...any) {
	s.Lock()
	defer s.Unlock()
	s.logroll.Rows = append(s.logroll.Rows, fmt.Sprintf("[ERRO](fg:red) "+format, args...))
	s.logroll.ScrollDown()
}

func infoln(args ...any) {
	if notui {
		logrus.Infoln(args...)
	} else {
		sc.infoln(args...)
	}
}

func (s *screen) infoln(args ...any) {
	s.Lock()
	defer s.Unlock()
	s.logroll.Rows = append(s.logroll.Rows, strings.TrimSuffix(
		"[INFO](fg:blue) "+fmt.Sprintln(args...), "\n",
	))
	s.logroll.ScrollDown()
}

func errorln(args ...any) {
	if notui {
		logrus.Errorln(args...)
	} else {
		sc.errorln(args...)
	}
}

func (s *screen) errorln(args ...any) {
	s.Lock()
	defer s.Unlock()
	s.logroll.Rows = append(s.logroll.Rows, strings.TrimSuffix(
		"[ERRO](fg:red) "+fmt.Sprintln(args...), "\n",
	))
	s.logroll.ScrollDown()
}

type progressmeter struct {
	prefix string
	name   string
	size   int
	prgs   int
	lstp   int
	fptr   *widgets.Gauge
	io.Writer
}

func newmeter(prefix, name string, size int) (pm progressmeter) {
	pm.prefix = prefix
	pm.name = name
	pm.size = size
	if !notui {
		pm.fptr = sc.addfile(prefix+" "+name, size)
	}
	return
}

func (pm *progressmeter) Write(p []byte) (n int, err error) {
	if pm.size == 0 {
		return 0, errZeroMeterSize
	}
	pm.prgs += len(p)
	if !notui {
		sc.logwrite(len(p))
	}
	percent := pm.prgs * 100 / pm.size
	if percent == pm.lstp {
		return len(p), nil
	}
	if notui {
		logrus.Infof("%s [%2d%%] %s\t(%dMB/%dMB)", pm.prefix, percent, pm.name, pm.prgs/1024/1024, pm.size/1024/1024)
	} else {
		pm.fptr.Percent = percent
		pm.fptr.Label = fmt.Sprintf("%d%% (%dMB/%dMB)", percent, pm.prgs/1024/1024, pm.size/1024/1024)
	}
	pm.lstp = percent
	return len(p), nil
}

func (pm *progressmeter) finish() {
	if !notui {
		sc.removefile(pm.fptr)
	}
}
