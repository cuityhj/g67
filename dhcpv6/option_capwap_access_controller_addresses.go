package dhcpv6

import (
	"fmt"
	"net"

	"github.com/u-root/uio/uio"
)

func OptCAPWAPAccessControllerAddresses(ips ...net.IP) Option {
	return &optCAPWAPAccessControllerAddresses{Addresses: ips}
}

type optCAPWAPAccessControllerAddresses struct {
	Addresses []net.IP
}

func (o *optCAPWAPAccessControllerAddresses) Code() OptionCode {
	return OptionCAPWAPAccessControllerAddresses
}

func (o *optCAPWAPAccessControllerAddresses) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	for _, ns := range o.Addresses {
		buf.WriteBytes(ns.To16())
	}

	return buf.Data()
}

func (o *optCAPWAPAccessControllerAddresses) String() string {
	return fmt.Sprintf("%s: %v", o.Code(), o.Addresses)
}

func (o *optCAPWAPAccessControllerAddresses) FromBytes(data []byte) error {
	buf := uio.NewBigEndianBuffer(data)
	o.Addresses = make([]net.IP, 0, buf.Len()/net.IPv6len)
	for buf.Has(net.IPv6len) {
		o.Addresses = append(o.Addresses, buf.CopyN(net.IPv6len))
	}

	return buf.FinError()
}
