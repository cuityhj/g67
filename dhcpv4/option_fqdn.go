package dhcpv4

import (
	"github.com/u-root/uio/uio"

	"github.com/cuityhj/g67/rfc1035label"
)

const (
	FLAG_S    = 0x01
	FLAG_O    = 0x02
	FLAG_E    = 0x04
	FLAG_N    = 0x08
	FLAG_MASK = 0xF
)

type FQDN struct {
	Flags       uint8
	Rcode1      uint8
	Rcode2      uint8
	DomainName  string
	domainBytes []byte
}

func (fqdn *FQDN) FromBytes(data []byte) (err error) {
	buf := uio.NewBigEndianBuffer(data)
	fqdn.Flags = buf.Read8()
	fqdn.Rcode1 = buf.Read8()
	fqdn.Rcode2 = buf.Read8()
	fqdn.domainBytes = buf.ReadAll()
	if fqdn.Flags&FLAG_E == FLAG_E {
		if domainName, err := rfc1035label.FromBytes(fqdn.domainBytes); err != nil {
			return err
		} else {
			fqdn.DomainName = domainName.String()
		}
	} else {
		fqdn.DomainName = string(fqdn.domainBytes)
	}

	return buf.FinError()
}

func (fqdn *FQDN) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write8(fqdn.Flags)
	buf.Write8(fqdn.Rcode1)
	buf.Write8(fqdn.Rcode2)
	buf.WriteBytes(fqdn.domainBytes)
	return buf.Data()
}

func (fqdn *FQDN) String() string {
	return fqdn.DomainName
}

func OptFQDN(fqdn *FQDN) Option {
	return Option{
		Code:  OptionFQDN,
		Value: fqdn,
	}
}
