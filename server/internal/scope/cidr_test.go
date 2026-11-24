package scope

import (
	"net/netip"
	"testing"
)

func TestParseCIDRRejectsHostBits(t *testing.T) {
	if _, err := ParseCIDR("10.20.30.5/24"); err == nil {
		t.Fatal("expected CIDR with host bits to be rejected")
	}
}

func TestContainsIPRespectsPrefixBoundary(t *testing.T) {
	prefix, err := ParseCIDR("10.20.30.0/24")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		address string
		want    bool
	}{
		{address: "10.20.30.1", want: true},
		{address: "10.20.31.1", want: false},
		{address: "10.20.30.255", want: true},
	}
	for _, test := range tests {
		if got := ContainsIP(prefix, netip.MustParseAddr(test.address)); got != test.want {
			t.Errorf("ContainsIP(%s) = %t, want %t", test.address, got, test.want)
		}
	}
}

func TestContainsIPv6AndMappedIPv4(t *testing.T) {
	v6Prefix, err := ParseCIDR("2001:db8:1::/48")
	if err != nil {
		t.Fatal(err)
	}
	if !ContainsIP(v6Prefix, netip.MustParseAddr("2001:db8:1::5")) {
		t.Fatal("expected IPv6 address inside prefix")
	}
	v4Prefix, err := ParseCIDR("192.0.2.0/24")
	if err != nil {
		t.Fatal(err)
	}
	if !ContainsIP(v4Prefix, netip.MustParseAddr("::ffff:192.0.2.7")) {
		t.Fatal("expected IPv4-mapped address inside IPv4 prefix")
	}
}
