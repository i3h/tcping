package main

import (
	"flag"
	"fmt"
	"os"
)

var Version bool
var PingDst string
var TCPingDst string
var HTTPingDst string
var MtrDst string
var Query string

var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

var Usage = func() {
	fmt.Fprintf(CommandLine.Output(), "Usage:\n")
	flag.PrintDefaults()
}

func init_flag() {
	flag.BoolVar(&Version, "v", false, "Version")
	flag.StringVar(&PingDst, "i", "", "ICMP Ping")
	flag.StringVar(&TCPingDst, "t", "", "TCP Ping")
	flag.StringVar(&HTTPingDst, "h", "", "HTTP Ping")
	flag.StringVar(&MtrDst, "m", "", "MTR Trace")
	flag.StringVar(&Query, "q", "", "Query ip information")
	flag.Parse()
}

func hasFlag() bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		found = true
	})
	return found
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
