package scope

import (
	"fmt"
	"net/netip"
)

func ParseCIDR(value string) (netip.Prefix, error) {
	prefix, err := netip.ParsePrefix(value)
	if err != nil {
		return netip.Prefix{}, fmt.Errorf("invalid CIDR %q: %w", value, err)
	}
	if prefix != prefix.Masked() {
		return netip.Prefix{}, fmt.Errorf("CIDR %q includes host bits; use %q", value, prefix.Masked())
	}
	return prefix, nil
}

func ContainsIP(prefix netip.Prefix, address netip.Addr) bool {
	if !prefix.IsValid() || !address.IsValid() {
		return false
	}
	return prefix.Contains(address.Unmap())
}
