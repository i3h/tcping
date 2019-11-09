package main

import (
	"flag"
	"fmt"
	"net"
	"strings"

	"github.com/noobly314/pingme/httping"
	"github.com/noobly314/pingme/mtr"
	"github.com/noobly314/pingme/ping"
	"github.com/noobly314/pingme/tcping"
)

var (
	VersionString string
)

func init() {
	init_log()
	init_flag()
}

func main() {
	if !hasFlag() {
		switch len(flag.Args()) {
		case 0:
			flag.PrintDefaults()
		case 1:
			addr := flag.Args()[0]

			// Query
			address := parseInput(addr)
			info := queryInfo(address)
			logQuery(info)

			// HTTP Ping
			if strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://") {
				fmt.Println()
				stats, err := httping.New(addr)
				logHttping(stats, err, addr)
			}
		case 2:
			addr := flag.Args()[0]
			port := flag.Args()[1]
			ip := lookupIP(addr)
			address := net.JoinHostPort(ip, port)
			c := tcping.New(address)
			logTcping(c, address)
		default:
			log.Warn("Too many arguments.")
		}
	} else {
		if isFlagPassed("v") {
			// Version
			fmt.Println(VersionString)
		} else if isFlagPassed("i") {
			// ICMP Ping
			dst, dur, err := ping.New(PingDst)
			logPing(dst, dur, err)
		} else if isFlagPassed("t") {
			// TCP Ping
			c := tcping.New(TCPingDst)
			logTcping(c, TCPingDst)
		} else if isFlagPassed("h") {
			// HTTP Ping
			stats, err := httping.New(HTTPingDst)
			logHttping(stats, err, HTTPingDst)
		} else if isFlagPassed("m") {
			// MTR
			hops, err := mtr.New(MtrDst)
			if err != nil {
				log.Fatal(err)
			}
			logMtr(hops, MtrDst)
		} else if isFlagPassed("q") {
			address := parseInput(Query)
			info := queryInfo(address)
			logQuery(info)
		}
	}
}
