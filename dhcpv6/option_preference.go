package dhcpv6

import (
	"fmt"

	"github.com/u-root/uio/uio"
)

func WithPreference(prefValue uint8) Modifier {
	return func(d DHCPv6) {
		d.AddOption(&OptPreference{PrefValue: prefValue})
	}
}

type OptPreference struct {
	PrefValue uint8
}

func (o *OptPreference) Code() OptionCode {
	return OptionPreference
}

func (o *OptPreference) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write8(o.PrefValue)
	return buf.Data()
}

func (o *OptPreference) String() string {
	return fmt.Sprintf("%s: %d", o.Code(), o.PrefValue)
}

func (o *OptPreference) FromBytes(data []byte) error {
	buf := uio.NewBigEndianBuffer(data)
	o.PrefValue = buf.Read8()
	return buf.FinError()
}
