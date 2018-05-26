package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"
)

var DEBUG_TARGETS []string
var rv *regexp.Regexp

func init() {
	rv = regexp.MustCompile(`.func\d+(.\d+)?\s*$`)
}

func Error(err error) {
	fmt.Fprintln(os.Stderr, trace(err.Error()))
}

func Log(msg string, args ...interface{}) {
	msg = strings.TrimRightFunc(fmt.Sprintf(msg, args...), unicode.IsSpace)
	fmt.Println(msg)
}

func Dbg(msg string, args ...interface{}) {
	if len(DEBUG_TARGETS) == 0 {
		return
	}
	var wanted bool
	caller := ""
	log := trace("")
	for _, l := range log {
		if l != "" {
			caller = l
			break
		}
	}
	caller = rv.ReplaceAllString(caller, "")
	if DEBUG_TARGETS[0] == "*" {
		wanted = true
	} else {
		if caller == "" {
			wanted = true
		} else {
			for _, t := range DEBUG_TARGETS {
				if strings.HasSuffix(caller, t) {
					wanted = true
					break
				}
			}
		}
	}
	if wanted {
		Log(strings.TrimSpace(caller)+"> "+msg, args...)
	}
}

func SetDebugTargets(targets string) {
	DEBUG_TARGETS = []string{}
	for _, t := range strings.Split(targets, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			DEBUG_TARGETS = append(DEBUG_TARGETS, t)
		}
	}
}

func Perf(tag string, work func()) {
	start := time.Now()
	Dbg("[EXEC]%s", tag)
	work()
	elapsed := time.Since(start).Seconds()
	Dbg("[DONE]%s (elapsed: %f)", tag, elapsed)
}
