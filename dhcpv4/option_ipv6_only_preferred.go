package dhcpv4

import (
	"strconv"

	"github.com/u-root/uio/uio"
)

type IPv6OnlyPreferred uint32

func (i IPv6OnlyPreferred) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write32(uint32(i))
	return buf.Data()
}

func (i IPv6OnlyPreferred) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

func (i *IPv6OnlyPreferred) FromBytes(data []byte) error {
	buf := uio.NewBigEndianBuffer(data)
	*i = IPv6OnlyPreferred(buf.Read32())
	return buf.FinError()
}

func OptIPv6OnlyPreferred(preferred uint32) Option {
	return Option{
		Code:  OptionIPv6OnlyPreferred,
		Value: IPv6OnlyPreferred(preferred),
	}
}
