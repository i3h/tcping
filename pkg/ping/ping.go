package ping

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type IPType struct {
	Type               string
	ListenAddr         string
	Network            string
	ICMPNetwork        string
	ProtocolNumber     int
	RequestMessageType icmp.Type
	ReplyMessageType   icmp.Type
}

var (
	IPType4 = IPType{
		Type:               "4",
		ListenAddr:         "0.0.0.0",
		Network:            "ip4",
		ICMPNetwork:        "ip4:icmp",
		ProtocolNumber:     1,
		RequestMessageType: ipv4.ICMPTypeEcho,
		ReplyMessageType:   ipv4.ICMPTypeEchoReply,
	}
	IPType6 = IPType{
		Type:               "6",
		ListenAddr:         "::",
		Network:            "ip6",
		ICMPNetwork:        "ip6:ipv6-icmp",
		ProtocolNumber:     58,
		RequestMessageType: ipv6.ICMPTypeEchoRequest,
		ReplyMessageType:   ipv6.ICMPTypeEchoReply,
	}
)

func New(address string) (*net.IPAddr, time.Duration, error) {
	// Check ip type
	// Resolve address
	var err error
	var dst *net.IPAddr
	var ipType IPType
	dst, err = net.ResolveIPAddr("ip4", address)
	if err != nil {
		dst, err = net.ResolveIPAddr("ip6", address)
		if err != nil {
			return nil, 0, err
		} else {
			ipType = IPType6
		}
	} else {
		ipType = IPType4
	}

	// Start listening for icmp replies
	c, err := icmp.ListenPacket(ipType.ICMPNetwork, ipType.ListenAddr)
	if err != nil {
		return nil, 0, err
	}
	defer c.Close()

	// Make a new ICMP message
	m := icmp.Message{
		Type: ipType.RequestMessageType,
		Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte(""),
		},
	}
	b, err := m.Marshal(nil)
	if err != nil {
		return dst, 0, err
	}

	// Send it
	start := time.Now()
	n, err := c.WriteTo(b, dst)
	if err != nil {
		return dst, 0, err
	} else if n != len(b) {
		return dst, 0, fmt.Errorf("got %v; want %v", n, len(b))
	}

	// Wait for a reply
	reply := make([]byte, 1500)
	err = c.SetReadDeadline(time.Now().Add(3 * time.Second))
	if err != nil {
		return dst, 0, err
	}
	n, peer, err := c.ReadFrom(reply)
	if err != nil {
		return dst, 0, err
	}
	duration := time.Since(start)

	// Pack it up boys, we're done here
	rm, err := icmp.ParseMessage(ipType.ProtocolNumber, reply[:n])
	if err != nil {
		return dst, 0, err
	}

	//return dst, duration, nil
	switch rm.Type {
	case ipType.ReplyMessageType:
		return dst, duration, nil
	default:
		return dst, 0, fmt.Errorf("got %+v from %v; want echo reply", rm, peer)
	}

	return dst, 0, err
}
