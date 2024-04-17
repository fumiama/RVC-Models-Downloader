package main

import (
	"strings"
)

type commandlist []string

var cmdlst = make(commandlist, 0, 64)

func (cl commandlist) String() string {
	sb := strings.Builder{}
	islastdir := false
	isfirstloop := true
	sb.WriteString("    ")
	for _, cmd := range cmdlst {
		if len(cmd) == 0 {
			continue
		}
		trimedcmd := strings.TrimSuffix(cmd, "/")
		a := strings.LastIndex(trimedcmd, "/") + 1
		b := len(cmd) - 1
		if a >= b {
			continue
		}
		isdir := cmd[b] == '/'
		ident := strings.Count(trimedcmd, "/") + 1
		if !isfirstloop && (islastdir || isdir) {
			sb.WriteByte('\n')
			if !isdir {
				ident--
			}
			for i := 0; i < ident; i++ {
				sb.WriteString("    ")
			}
		}
		isfirstloop = false
		if isdir {
			islastdir = true
			sb.WriteString(cmd[a:b])
			sb.WriteByte(':')
		} else { // is file
			islastdir = false
			sb.WriteString("    ")
			sb.WriteString(strings.TrimSuffix(cmd[a:], ".yaml"))
		}
	}
	return sb.String()
}
