package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
	"unsafe"

	"github.com/fumiama/terasu/ip"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type screen struct {
	sync.RWMutex
	sysinfo  *widgets.Paragraph
	logroll  *widgets.List
	speedln  *widgets.Plot
	prgbars  []ui.Drawable
	reusepg  []*widgets.Gauge
	currh, w int
	totaldl  int
	lastclr  time.Time
}

func newscreen() (s screen) {
	w, h := ui.TerminalDimensions()
	s.w = w

	s.sysinfo = widgets.NewParagraph()
	s.sysinfo.Title = "System Info"
	s.sysinfo.BorderStyle.Fg = ui.ColorGreen
	s.sysinfo.Text = fmt.Sprintf(
		"[Time](mod:bold): %s\n[OS](mod:bold): %s, [Architecture](mod:bold): %s\n[Is IPv6 available](mod:bold): %v",
		time.Now().Local().Format(time.DateTime+" (MST)"),
		runtime.GOOS, runtime.GOARCH,
		ip.IsIPv6Available.Get(),
	)
	s.sysinfo.SetRect(0, s.currh, w/2, s.currh+5)
	s.currh += 5

	s.logroll = widgets.NewList()
	s.logroll.Title = "Logs"
	s.logroll.BorderStyle.Fg = ui.ColorBlue
	s.logroll.WrapText = false
	s.logroll.SetRect(w/2, 0, w, h/2)

	s.speedln = widgets.NewPlot()
	s.speedln.Title = "Speed"
	s.speedln.Data = make([][]float64, 1)
	s.speedln.Data[0] = []float64{0, 0}
	s.speedln.AxesColor = ui.ColorWhite
	s.speedln.LineColors[0] = ui.ColorYellow
	s.speedln.BorderStyle.Fg = ui.ColorYellow
	s.speedln.SetRect(w/2, h/2, w, h)
	return
}

func (s *screen) logwrite(sz int) {
	s.Lock()
	defer s.Unlock()
	s.totaldl += sz
	tdiff := time.Since(s.lastclr)
	if tdiff > time.Second {
		s.speedln.Data[0] = append(s.speedln.Data[0],
			float64(s.totaldl/1024)/(float64(tdiff)/float64(time.Second)),
		)
		s.totaldl = 0
		s.lastclr = time.Now()
	}
}

func (s *screen) flushloop(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	s.flush()
	for {
		select {
		case e := <-ui.PollEvents():
			s.flush()
			if e.Type == ui.KeyboardEvent {
				switch e.ID {
				case "q", "<Escape>", "<C-c>":
					return
				}
			}
		case <-t.C:
			s.flush()
		}
	}
}

func (s *screen) flush() {
	s.RLock()
	defer s.RUnlock()
	ui.Render(s.sysinfo, s.logroll, s.speedln)
	if len(s.prgbars) > 0 {
		ui.Render(s.prgbars...)
	}
}

func (s *screen) addfile(name string, size int) *widgets.Gauge {
	s.Lock()
	defer s.Unlock()
	var g *widgets.Gauge
	if len(s.reusepg) > 0 {
		b := len(s.reusepg) - 1
		g = s.reusepg[b]
		s.reusepg = s.reusepg[:b]
	} else {
		g = widgets.NewGauge()
		g.SetRect(0, s.currh, s.w/2, s.currh+3)
		s.currh += 3
	}
	g.Title = name
	g.Label = fmt.Sprintf("0%% (0MB/%dMB)", size/1024/1024)
	s.prgbars = append(s.prgbars, g)
	return g
}

func (s *screen) removefile(g *widgets.Gauge) {
	s.Lock()
	defer s.Unlock()
	for i, obj := range s.prgbars {
		if *(**widgets.Gauge)(unsafe.Add(
			unsafe.Pointer(&obj), unsafe.Sizeof(uintptr(0)),
		)) == g {
			switch i {
			case 0:
				s.prgbars = s.prgbars[1:]
			case len(s.prgbars) - 1:
				s.prgbars = s.prgbars[:len(s.prgbars)-1]
			default:
				s.prgbars = append(s.prgbars[:i], s.prgbars[i+1:]...)
			}
			s.reusepg = append(s.reusepg, g)
			return
		}
	}
}
