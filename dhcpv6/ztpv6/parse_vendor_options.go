package ztpv6

import (
	"errors"
	"strconv"
	"strings"

	"github.com/cuityhj/g67/dhcpv6"
	"github.com/cuityhj/g67/iana"
)

var (
	errVendorOptionMalformed = errors.New("malformed vendor option")
)

// VendorData contains fields extracted from Option 17 data
type VendorData struct {
	VendorName, Model, Serial string
}

// ParseVendorData will try to parse dhcp6 Vendor Specific Information options
// ( 16 and 17) data looking for more specific vendor data (like model, serial
// number, etc). If the options are missing we will just return nil
func ParseVendorData(packet dhcpv6.DHCPv6) (*VendorData, error) {
	// check for both options 16 and 17 if both are present will use opt 17
	opt16 := packet.GetOneOption(dhcpv6.OptionVendorClass)
	opt17 := packet.GetOneOption(dhcpv6.OptionVendorOpts)
	if (opt16 == nil) && (opt17 == nil) {
		return nil, errors.New("no vendor options or vendor class found")
	}

	vd := VendorData{}
	vData := []string{}

	if opt17 != nil {
		vendorOptsOption := opt17.(*dhcpv6.OptVendorOpts)

		// MLNX-OS has the relevant information spread over different sub-options
		// of option 17 so the usual approach doesn't work
		if vendorOptsOption.EnterpriseNumber == uint32(iana.EnterpriseIDMellanoxTechnologiesLTD) {
			return getMellanoxVendorData(vendorOptsOption)
		}

		// rest of vendors use a single sub-option so we stringify them and parse them below
		for _, opt := range vendorOptsOption.VendorOpts {
			vData = append(vData, string(opt.(*dhcpv6.OptionGeneric).OptionData))
		}
	} else {
		data := opt16.(*dhcpv6.OptVendorClass).Data
		for _, d := range data {
			vData = append(vData, string(d))
		}
	}
	for _, d := range vData {
		switch {
		// Arista;DCS-0000;00.00;ZZZ00000000
		// Cisco;8800;12.34;FOC00000000
		case strings.HasPrefix(d, "Arista;"), strings.HasPrefix(d, "Cisco;"):
			p := strings.Split(d, ";")
			if len(p) < 4 {
				return nil, errVendorOptionMalformed
			}

			vd.VendorName = p[0]
			vd.Model = p[1]
			vd.Serial = p[3]
			return &vd, nil

		// ZPESystems:NSC:000000000
		case strings.HasPrefix(d, "ZPESystems:"):
			p := strings.Split(d, ":")
			if len(p) < 3 {
				return nil, errVendorOptionMalformed
			}

			vd.VendorName = p[0]
			vd.Model = p[1]
			vd.Serial = p[2]
			return &vd, nil
		
		// NVOS##MMM1234##MM1234X56ABC
		case strings.HasPrefix(d, "NVOS##"):
			p := strings.Split(d, "##")
			if len(p) < 3 {
				return nil, errVendorOptionMalformed
			}

			vd.VendorName = p[0]
			vd.Model = p[1]
			vd.Serial = p[2]
			return &vd, nil

		// For Ciena the class identifier (opt 60) is written in the following format:
		//    {vendor iana code}-{product}-{type}
		// For Ciena the iana code is 1271
		// The product type is a number that maps to a Ciena product
		// The type is used to identified different subtype of the product.
		// An example can be ‘1271-23422Z11-123’.
		case strings.HasPrefix(d, strconv.Itoa(int(iana.EnterpriseIDCienaCorporation))):
			v := strings.Split(d, "-")
			if len(v) < 3 {
				return nil, errVendorOptionMalformed
			}
			vd.VendorName = iana.EnterpriseIDCienaCorporation.String()
			vd.Model = v[1] + "-" + v[2]
			duid := packet.(*dhcpv6.Message).Options.ClientID()
			if enterpriseDUID, ok := duid.(*dhcpv6.DUIDEN); ok {
				vd.Serial = string(enterpriseDUID.EnterpriseIdentifier)
			}
			return &vd, nil
		}
	}
	return nil, errors.New("failed to parse vendor option data")
}
