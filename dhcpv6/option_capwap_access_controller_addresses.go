package dhcpv6

import (
	"net"

	"github.com/u-root/uio/uio"
)

func OptCapwapAccessControllerAddresses(ips ...net.IP) Option {
	return &OptionCapwapAccessControllerAddresses{Addresses: ips}
}

type OptionCapwapAccessControllerAddresses struct {
	Addresses []net.IP
}

func (o *OptionCapwapAccessControllerAddresses) Code() OptionCode {
	return OptionCAPWAPAccessControllerAddresses
}

func (o *OptionCapwapAccessControllerAddresses) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	for _, ns := range o.Addresses {
		buf.WriteBytes(ns.To16())
	}

	return buf.Data()
}

func (o *OptionCapwapAccessControllerAddresses) String() string {
	return fmt.Sprintf("%v", o.Addresses)
}

func (o *OptionCapwapAccessControllerAddresses) FromBytes(data []byte) error {
	buf := uio.NewBigEndianBuffer(data)
	for buf.Has(net.IPv6len) {
		o.Addresses = append(o.Addresses, buf.CopyN(net.IPv6len))
	}

	return buf.FinError()
}
