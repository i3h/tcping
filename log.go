package main

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/noobly314/pingme/httping"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

var (
	cyan   = color.New(color.FgCyan).SprintfFunc()
	blue   = color.New(color.FgBlue).SprintfFunc()
	green  = color.New(color.FgGreen).SprintfFunc()
	yellow = color.New(color.FgYellow).SprintfFunc()
	red    = color.New(color.FgRed).SprintfFunc()
)

func init_log() {
	formatter := &logrus.TextFormatter{
		DisableTimestamp: true,
	}
	log.SetFormatter(formatter)
	log.Out = os.Stdout
}

func logHttping(stats httping.Stats, err error, address string) {
	if err == nil {
		fmt.Printf("%s:    %s\n", cyan("%-10s", "Proxy"), strconv.FormatBool(stats.Proxy))
		fmt.Printf("%s:    %s\n", cyan("%-10s", "Scheme"), stats.Scheme)
		fmt.Printf("%s:    %s\n", cyan("%-10s", "Host"), parseInput(address))
		fmt.Printf("%s:    %.2f ms\n", cyan("%-10s", "DNS Lookup"), float64(stats.DNS)/1e6)
		fmt.Printf("%s:    %.2f ms\n", cyan("%-10s", "TCP"), float64(stats.TCP)/1e6)
		if stats.Scheme == "https" {
			fmt.Printf("%s:    %.2f ms\n", cyan("%-10s", "TLS"), float64(stats.TLS)/1e6)
		}
		fmt.Printf("%s:    %.2f ms\n", cyan("%-10s", "Process"), float64(stats.Process)/1e6)
		fmt.Printf("%s:    %.2f ms\n", cyan("%-10s", "Transfer"), float64(stats.Transfer)/1e6)
		fmt.Printf("%s:    %.2f ms\n", cyan("%-10s", "Total"), float64(stats.Total)/1e6)
	}
}

func logTcping(code int, address string) {
	if code == 0 {
		fmt.Printf("%s%s%s\n", cyan("%-7s", "TCP"), green("%-10s", "OPEN"), address)
	} else if code == 1 {
		fmt.Printf("%s%s%s\n", cyan("%-7s", "TCP"), yellow("%-10s", "CLOSED"), address)
	} else if code == 2 {
		fmt.Printf("%s%s%s\n", cyan("%-7s", "TCP"), red("%-10s", "ERROR"), address)
	}
}

func logPing(dst *net.IPAddr, dur time.Duration, err error) {
	if err != nil {
		match, _ := regexp.MatchString("operation not permitted", err.Error())
		if match {
			fmt.Printf("%s%s%s\n", cyan("%-7s", "ICMP"), red("%-10s", "ERROR"), red("No privileges"))
		} else {
			fmt.Printf("%s%s%s\n", cyan("%-7s", "ICMP"), red("%-10s", "ERROR"), dst.String())
		}
		return
	}
	fmt.Printf("%s%s%s    %s ms\n", cyan("%-7s", "ICMP"), green("%-10s", "OPEN"), dst.String(), fmt.Sprintf("%.1f", float64(dur.Microseconds())/1000))
}

func logMtr(hops []string, address string) {
	for _, h := range hops {
		fmt.Printf("%s%s\n", cyan("%-7s", "MTR"), h)
	}
}

func logQuery(info IPInfo) {
	v := reflect.ValueOf(info)
	names := make([]string, v.NumField())
	values := make([]string, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		names[i] = v.Type().Field(i).Name
		values[i] = v.Field(i).Interface().(string)
	}
	l := getMaxNameLength(names)
	for i := 0; i < v.NumField(); i++ {
		fmt.Printf("%s:    %s\n", cyan("%-*s", l, names[i]), values[i])
	}
}

func getMaxNameLength(names []string) int {
	var length int
	for _, val := range names {
		if len(val) > length {
			length = len(val)
		}
	}
	return length
}
