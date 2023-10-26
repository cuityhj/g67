package dhcpv4

import (
	"strings"

	"github.com/u-root/uio/uio"

	"github.com/cuityhj/g67/rfc1035label"
)

type FQDN struct {
	Flags      uint8
	Rcode1     uint8
	Rcode2     uint8
	DomainName *rfc1035label.Labels
}

func (fqdn *FQDN) FromBytes(data []byte) (err error) {
	buf := uio.NewBigEndianBuffer(data)
	fqdn.Flags = buf.Read8()
	fqdn.Rcode1 = buf.Read8()
	fqdn.Rcode2 = buf.Read8()
	if fqdn.DomainName, err = rfc1035label.FromBytes(buf.ReadAll()); err != nil {
		return err
	} else {
		return buf.FinError()
	}
}

func (fqdn FQDN) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write8(fqdn.Flags)
	buf.Write8(fqdn.Rcode1)
	buf.Write8(fqdn.Rcode2)
	buf.WriteBytes(fqdn.DomainName.ToBytes())
	return buf.Data()
}

func (fqdn FQDN) String() string {
	if fqdn.DomainName == nil || fqdn.DomainName.String() == "[]" {
		return ""
	} else {
		return strings.TrimSuffix(strings.TrimPrefix(fqdn.DomainName.String(), "["), "]")
	}
}

func OptFQDN(fqdn *FQDN) Option {
	return Option{
		Code:  OptionFQDN,
		Value: fqdn,
	}
}
