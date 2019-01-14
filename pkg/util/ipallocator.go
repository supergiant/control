package util

import (
	"math/big"
	"net"

	"github.com/pkg/errors"
)

var (
	ErrSubnetTooSmall = errors.New("subnet too smal")
)

// TODO: move to the separate package?
// GetDNSIP returns a dnsIP, which is 10th IP in svcSubnet CIDR range
func GetKubernetesDefaultSvcIP(svcSubnet string) (net.IP, error) {
	return indexedIpFrom(svcSubnet, 1)
}

// GetDNSIP returns a dnsIP, which is 10th IP in svcSubnet CIDR range
func GetDNSIP(svcSubnet string) (net.IP, error) {
	return indexedIpFrom(svcSubnet, 10)
}

func indexedIpFrom(subnet string, index int) (net.IP, error) {
	// Get the service subnet CIDR
	_, svcSubnetCIDR, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't parse service subnet")
	}

	// Selects the 10th IP in service subnet CIDR range as dnsIP
	return GetIndexedIP(svcSubnetCIDR, index)
}

// GetIndexedIP returns a net.IP that is subnet.IP + index in the contiguous IP space.
// Has been copied from the source to reduce number of dependencies.
// src: https://github.com/kubernetes/kubernetes/blob/master/pkg/registry/core/service/ipallocator/allocator.go#L280
func GetIndexedIP(subnet *net.IPNet, index int) (net.IP, error) {
	ip := addIPOffset(bigForIP(subnet.IP), index)
	if !subnet.Contains(ip) {
		return nil, errors.Wrapf(ErrSubnetTooSmall, "can't generate IP with index %d from %s subnet", index, subnet)
	}
	return ip, nil
}

// bigForIP creates a big.Int based on the provided net.IP
func bigForIP(ip net.IP) *big.Int {
	b := ip.To4()
	if b == nil {
		b = ip.To16()
	}
	return big.NewInt(0).SetBytes(b)
}

// addIPOffset adds the provided integer offset to a base big.Int representing a
// net.IP
func addIPOffset(base *big.Int, offset int) net.IP {
	return net.IP(big.NewInt(0).Add(base, big.NewInt(int64(offset))).Bytes())
}
