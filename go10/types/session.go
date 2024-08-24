package types

import "net"

type Session struct {
	Datagram *Datagram // The datagram associated with this session
	Addr     *net.UDPAdd
}
