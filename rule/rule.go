// Package rule  represents all the logic for the listener rules
package rule

import (
	"errors"
	"fmt"
)

// Protocol represents the valid protocols for an EC2 port listener
type Protocol int

// String implements stringer interface for RuleProtocol
func (r Protocol) String() string {
	switch r {
	case ALL:
		return "all"
	case TCP:
		return "tcp"
	case UDP:
		return "udp"
	case ICMP:
		return "icmp"
	default:
		return "Unknown"
	}
}

// ParseProtocol parses a valid protocol
func ParseProtocol(proto string) (Protocol, error) {
	switch proto {
	case "all", "ALL":
		return ALL, nil
	case "tcp", "TCP":
		return TCP, nil
	case "udp", "UDP":
		return UDP, nil
	case "icmp", "ICMP":
		return ICMP, nil
	default:
		return 0, errors.New("not a valid protocol")
	}
}

const (
	// ALL protocols
	ALL Protocol = iota
	// TCP protocol
	TCP
	// UDP Protocol
	UDP
	// ICMP Protocol
	ICMP
)

const (
	// AllCIDR is the CIDR that describes all addresses
	AllCIDR = "0.0.0.0/0"
)

// Rule defines the rule where the instances need to be accessed
type Rule struct {
	CIDR     string
	Port     int
	Protocol Protocol
}

// New creates a regular Rule
func New(cidr string, port int, protocol Protocol) *Rule {
	r := &Rule{
		CIDR:     cidr,
		Port:     port,
		Protocol: protocol,
	}

	return r
}

// NewOpenRule creates a rule to a port accessible with all supported protocols and from anywhere
func NewOpenRule(port int) *Rule {
	r := &Rule{
		CIDR:     AllCIDR,
		Port:     port,
		Protocol: ALL,
	}

	return r
}

// Implement stringer interface
func (r *Rule) String() string {
	return fmt.Sprintf("%s:%d:%s", r.CIDR, r.Port, r.Protocol)
}
