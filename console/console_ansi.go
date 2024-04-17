//go:build !windows

// Package console sets console's behavior on init
package console

// https://github.com/FloatTech/ZeroBot-Plugin/blob/master/console/console_ansi.go

import (
	"fmt"
)

func init() {
	fmt.Print("\033]0;RVC Models Downloader\007")
}
