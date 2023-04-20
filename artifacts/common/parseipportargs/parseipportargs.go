package parseipportargs

import (
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"strings"
)

func ValidatePortRange(p int) error {
	if p < 0 || p > 65535 {
		return fmt.Errorf("port out of range 0-65535: %v", p)
	}
	return nil
}

func aToPort(s string) (int, error) {
	p, err := strconv.Atoi(s)
	if err != nil {
		return -1, err
	}
	err = ValidatePortRange(p)
	if err != nil {
		return -1, err
	}
	return p, nil
}

func uniqIP(IPs []netip.Addr) ([]netip.Addr, error) {
	set := make(map[netip.Addr]bool)
	uniqIPs := make([]netip.Addr, 0)

	for _, IP := range IPs {
		// netip.Addr is comparable, so that it can be used
		// as a key.  Note that, though string is also
		// comparable, it cannot properly compare different
		// IPv6 strings that map to the same IP.
		if set[IP] {
			continue
		}
		set[IP] = true
		uniqIPs = append(uniqIPs, IP)
	}
	return uniqIPs, nil
}

func uniqPort(ports []int) []int {
	set := make(map[int]bool)
	uniqPorts := make([]int, 0)
	for _, p := range ports {
		if set[p] {
			continue
		}
		set[p] = true
		uniqPorts = append(uniqPorts, p)
	}
	return uniqPorts
}

func ExpandCIDR(cidr string) ([]netip.Addr, error) {
	ips := make([]netip.Addr, 0)

	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return nil, err
	}
	for addr := prefix.Masked().Addr(); prefix.Contains(addr); addr = addr.Next() {
		ips = append(ips, addr)
	}
	return ips, nil
}

func NetipToNet(ips []netip.Addr) ([]net.IP, error) {
	IPs := make([]net.IP, 0)
	for _, ip := range ips {
		IP := net.ParseIP(ip.String())
		if IP == nil {
			return nil, fmt.Errorf("Error happened when converting %v to net.IP", IP)
		}
		IPs = append(IPs, IP)
	}
	return IPs, nil
}

func ParseIPArgs(s string) ([]net.IP, error) {
	// use netip.Addr as it is comparable
	IPs := make([]netip.Addr, 0)
	for _, ipOrCidr := range strings.Split(s, ",") {
		if strings.Contains(ipOrCidr, "/") {
			ips, err := ExpandCIDR(ipOrCidr)
			if err != nil {
				return nil, fmt.Errorf("Error happened when parsing CIDR %+q: %s", ipOrCidr, err)
			}
			IPs = append(IPs, ips...)
		} else {
			ip, err := netip.ParseAddr(ipOrCidr)
			if err != nil {
				return nil, fmt.Errorf("Invalid IP %+q: %s", ipOrCidr, err)
			}
			IPs = append(IPs, ip)
		}
	}
	// remove duplicates
	uniqIPs, err := uniqIP(IPs)
	if err != nil {
		return nil, err
	}
	// convert netip.Addr to more common net.IP
	netIPs, err := NetipToNet(uniqIPs)
	if err != nil {
		return nil, err
	}
	return netIPs, nil
}

func ParsePortArgs(s string) ([]int, error) {
	ports := make([]int, 0)
	for _, b := range strings.Split(s, ",") {
		k := strings.Split(b, "-")
		if len(k) == 1 {
			p, err := aToPort(k[0])
			if err != nil {
				return nil, err
			}
			ports = append(ports, p)
		} else if len(k) == 2 {
			low, err := aToPort(k[0])
			if err != nil {
				return nil, err
			}
			high, err := aToPort(k[1])
			if err != nil {
				return nil, err
			}
			if low > high {
				return nil, fmt.Errorf("port %v is higher than %v: %v", low, high, b)
			}
			for p := low; p <= high; p++ {
				ports = append(ports, p)
			}
		} else {
			return nil, fmt.Errorf("Invalid range syntax: %+q", b)
		}
	}
	// remove duplicates
	uniqPorts := uniqPort(ports)

	return uniqPorts, nil
}

// https://tools.ietf.org/html/rfc1035#section-3.2.2
// https://github.com/miekg/dns/blob/master/types.go#L24
var MapRRType = map[string]uint16{
	"None":       0,
	"A":          1,
	"NS":         2,
	"MD":         3,
	"MF":         4,
	"CNAME":      5,
	"SOA":        6,
	"MB":         7,
	"MG":         8,
	"MR":         9,
	"NULL":       10,
	"PTR":        12,
	"HINFO":      13,
	"MINFO":      14,
	"MX":         15,
	"TXT":        16,
	"RP":         17,
	"AFSDB":      18,
	"X25":        19,
	"ISDN":       20,
	"RT":         21,
	"NSAPPTR":    23,
	"SIG":        24,
	"KEY":        25,
	"PX":         26,
	"GPOS":       27,
	"AAAA":       28,
	"LOC":        29,
	"NXT":        30,
	"EID":        31,
	"NIMLOC":     32,
	"SRV":        33,
	"ATMA":       34,
	"NAPTR":      35,
	"KX":         36,
	"CERT":       37,
	"DNAME":      39,
	"OPT":        41,
	"APL":        42,
	"DS":         43,
	"SSHFP":      44,
	"RRSIG":      46,
	"NSEC":       47,
	"DNSKEY":     48,
	"DHCID":      49,
	"NSEC3":      50,
	"NSEC3PARAM": 51,
	"TLSA":       52,
	"SMIMEA":     53,
	"HIP":        55,
	"NINFO":      56,
	"RKEY":       57,
	"TALINK":     58,
	"CDS":        59,
	"CDNSKEY":    60,
	"OPENPGPKEY": 61,
	"CSYNC":      62,
	"ZONEMD":     63,
	"SVCB":       64,
	"HTTPS":      65,
	"SPF":        99,
	"UINFO":      100,
	"UID":        101,
	"GID":        102,
	"UNSPEC":     103,
	"NID":        104,
	"L32":        105,
	"L64":        106,
	"LP":         107,
	"EUI48":      108,
	"EUI64":      109,
	"URI":        256,
	"CAA":        257,
	"AVC":        258,

	"TKEY": 249,
	"TSIG": 250,

	// valid Question.Qtype only
	"IXFR":  251,
	"AXFR":  252,
	"MAILB": 253,
	"MAILA": 254,
	"ANY":   255,

	"TA":       32768,
	"DLV":      32769,
	"Reserved": 65535,
}

func ValidateRRTypeRange(p int) error {
	println(p)
	if p < 0 || p > 65535 {
		return fmt.Errorf("RRType out of range 0-65535: %v", p)
	}
	return nil
}

func aToRRType(RRType string) (uint16, error) {
	var val uint16
	valInt, err := strconv.Atoi(RRType)
	if err == nil {
		err = ValidateRRTypeRange(valInt)
		if err != nil {
			return val, err
		}
		// return value if it's already a value
		return uint16(valInt), nil
	}

	val, ok := MapRRType[RRType]
	if ok == false {
		return val, fmt.Errorf("Invalid RRType: %v", RRType)
	}
	return val, nil
}

func uniqRRType(RRTypes []uint16) []uint16 {
	set := make(map[uint16]bool)
	uniqRRTypes := make([]uint16, 0)
	for _, p := range RRTypes {
		if set[p] {
			continue
		}
		set[p] = true
		uniqRRTypes = append(uniqRRTypes, p)
	}
	return uniqRRTypes
}

func ParseRRTypeArgs(s string) ([]uint16, error) {
	RRTypes := make([]uint16, 0)
	for _, b := range strings.Split(s, ",") {
		k := strings.Split(b, "-")
		if len(k) == 1 {
			p, err := aToRRType(k[0])
			if err != nil {
				return nil, err
			}
			RRTypes = append(RRTypes, p)
		} else if len(k) == 2 {
			low, err := aToRRType(k[0])
			if err != nil {
				return nil, err
			}
			high, err := aToRRType(k[1])
			if err != nil {
				return nil, err
			}
			if low > high {
				return nil, fmt.Errorf("RRType %v is higher than %v: %v", low, high, b)
			}
			// important to cast to unit32 to avoid overflow when 65535++.
			for p := uint32(low); p <= uint32(high); p++ {
				RRTypes = append(RRTypes, uint16(p))
			}
		} else {
			return nil, fmt.Errorf("Invalid range syntax: %+q", b)
		}
	}
	// remove duplicates
	uniqRRTypes := uniqRRType(RRTypes)

	return uniqRRTypes, nil
}
