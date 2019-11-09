package tcping

import (
	"net"
	"regexp"
	"time"
)

func New(address string) int {
	d := net.Dialer{Timeout: 3 * time.Second}
	_, err := d.Dial("tcp", address)
	if err != nil {
		match, _ := regexp.MatchString("refused", err.Error())
		if match {
			// Closed
			return 1
		}
		match, _ = regexp.MatchString("timeout", err.Error())
		if match {
			// Timeout
			return 2
		}
	} else {
		// Open
		return 0
	}
	// Default
	return 2
}
