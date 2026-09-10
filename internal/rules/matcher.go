package rules

import (
	"fmt"
	"net"
	"strings"
)

// RuleType represents the type of rule
type RuleType int

const (
	RuleIPCIDR RuleType = iota
	RuleDomain
	RuleDomainSuffix
	RuleDomainKeyword
	RuleGeoSite
	RuleGeoIP
	RuleDefault
)

// Rule represents a routing rule
type Rule struct {
	Type    RuleType
	Value   string
	Network string // outbound network tag
}

// Matcher matches rules against a destination
type Matcher struct {
	rules []Rule
}

// NewMatcher creates a new rule matcher
func NewMatcher() *Matcher {
	return &Matcher{}
}

// AddRule adds a rule to the matcher
func (m *Matcher) AddRule(rule Rule) {
	m.rules = append(m.rules, rule)
}

// Match tries to match a destination against the rules
// Returns the network tag or empty string if no match
func (m *Matcher) Match(dest string) string {
	for _, rule := range m.rules {
		switch rule.Type {
		case RuleDomain:
			if dest == rule.Value {
				return rule.Network
			}
		case RuleDomainSuffix:
			if strings.HasSuffix(dest, rule.Value) {
				return rule.Network
			}
		case RuleDomainKeyword:
			if strings.Contains(dest, rule.Value) {
				return rule.Network
			}
		case RuleIPCIDR:
			if m.matchIPCIDR(dest, rule.Value) {
				return rule.Network
			}
		case RuleDefault:
			return rule.Network
		}
	}
	return ""
}

func (m *Matcher) matchIPCIDR(dest, cidr string) bool {
	ip := net.ParseIP(dest)
	if ip == nil {
		return false
	}

	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}

	return ipNet.Contains(ip)
}

// ParseRule parses a rule string into a Rule
func ParseRule(ruleStr string) (Rule, error) {
	parts := strings.SplitN(ruleStr, ",", 2)
	if len(parts) != 2 {
		return Rule{}, fmt.Errorf("invalid rule format: %s", ruleStr)
	}

	ruleType := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	var rule Rule
	switch strings.ToLower(ruleType) {
	case "ipcidr", "ip-cidr":
		rule.Type = RuleIPCIDR
	case "domain":
		rule.Type = RuleDomain
	case "domainsuffix", "domain-suffix":
		rule.Type = RuleDomainSuffix
	case "domainkeyword", "domain-keyword":
		rule.Type = RuleDomainKeyword
	case "geosite":
		rule.Type = RuleGeoSite
	case "geoip":
		rule.Type = RuleGeoIP
	default:
		return Rule{}, fmt.Errorf("unknown rule type: %s", ruleType)
	}

	rule.Value = value
	rule.Network = "proxy"
	return rule, nil
}
