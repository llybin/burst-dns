package dns

import (
	"errors"
	"golang.org/x/net/dns/dnsmessage"
	"net"
)

var (
	errIPInvalid   = errors.New("invalid IP address")
	errUnknownType = errors.New("type unknown")
)

func toResource(req request) (dnsmessage.Resource, error) {
	rName, err := dnsmessage.NewName(req.Host)
	none := dnsmessage.Resource{}
	if err != nil {
		return none, err
	}

	var rType dnsmessage.Type
	var rBody dnsmessage.ResourceBody

	switch req.Type {
	case "a":
		rType = dnsmessage.TypeA
		ip := net.ParseIP(req.Data)
		if ip == nil {
			return none, errIPInvalid
		}
		rBody = &dnsmessage.AResource{A: [4]byte{ip[12], ip[13], ip[14], ip[15]}}
	case "aaaa":
		rType = dnsmessage.TypeAAAA
		ip := net.ParseIP(req.Data)
		if ip == nil {
			return none, errIPInvalid
		}
		var ipV6 [16]byte
		copy(ipV6[:], ip)
		rBody = &dnsmessage.AAAAResource{AAAA: ipV6}
	default:
		return none, errUnknownType
	}

	return dnsmessage.Resource{
		Header: dnsmessage.ResourceHeader{
			Class: dnsmessage.ClassINET,
			Name:  rName,
			Type:  rType,
			TTL:   req.TTL,
		},
		Body: rBody,
	}, nil
}

func dnsTypeToStr(dnsType dnsmessage.Type) string {
	switch dnsType {
	case dnsmessage.TypeA:
		return "a"
	case dnsmessage.TypeAAAA:
		return "aaaa"
	default:
		return ""
	}
}
