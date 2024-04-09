package protocol

import (
	"context"
	"net"
	"strconv"
)

func resolveUDPAddr(ctx context.Context, addr string) (*net.UDPAddr, bool) {
	host, portString, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, false
	}
	portNumber, err := strconv.ParseUint(portString, 10, 16)
	if err != nil {
		return nil, false
	}
	if portNumber <= 0 {
		return nil, false
	}
	ipaddr, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, false
	}
	if len(ipaddr) <= 0 {
		return nil, false
	}
	return &net.UDPAddr{
		IP:   ipaddr[0].IP,
		Port: int(portNumber),
		Zone: ipaddr[0].Zone,
	}, true
}
