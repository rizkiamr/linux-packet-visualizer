package contract

// NetfilterHook represents a netfilter hook point where iptables/nftables
// rules are evaluated. These are the integration points for packet filtering,
// NAT, and packet mangling.
type NetfilterHook struct {
	// Hook is the netfilter hook name: PREROUTING, INPUT, FORWARD, OUTPUT, POSTROUTING
	Hook string `json:"hook"`

	// Tables lists the iptables tables traversed at this hook, in order
	// Possible tables: raw, mangle, nat, filter
	Tables []string `json:"tables"`

	// Description explains what happens at this hook point
	Description string `json:"description"`

	// Priority indicates the hook priority (lower = earlier)
	Priority int `json:"priority,omitempty"`
}

// Netfilter hook constants
const (
	HookPrerouting  = "PREROUTING"
	HookInput       = "INPUT"
	HookForward     = "FORWARD"
	HookOutput      = "OUTPUT"
	HookPostrouting = "POSTROUTING"
)

// NewOutputHook creates a netfilter OUTPUT hook annotation.
// OUTPUT is called for locally generated packets before routing.
func NewOutputHook() *NetfilterHook {
	return &NetfilterHook{
		Hook:        HookOutput,
		Tables:      []string{"raw", "mangle", "nat", "filter"},
		Description: "Locally generated packets. Firewall rules (iptables -A OUTPUT) are evaluated here.",
		Priority:    -100,
	}
}

// NewPostroutingHook creates a netfilter POSTROUTING hook annotation.
// POSTROUTING is called after routing, just before the packet leaves.
func NewPostroutingHook() *NetfilterHook {
	return &NetfilterHook{
		Hook:        HookPostrouting,
		Tables:      []string{"mangle", "nat"},
		Description: "Final hook before packet leaves. SNAT/MASQUERADE applied here.",
		Priority:    100,
	}
}

// NewPreroutingHook creates a netfilter PREROUTING hook annotation.
// PREROUTING is called for incoming packets before routing decision.
func NewPreroutingHook() *NetfilterHook {
	return &NetfilterHook{
		Hook:        HookPrerouting,
		Tables:      []string{"raw", "mangle", "nat"},
		Description: "First hook for incoming packets. DNAT applied here before routing.",
		Priority:    -300,
	}
}

// NewInputHook creates a netfilter INPUT hook annotation.
// INPUT is called for packets destined for the local machine.
func NewInputHook() *NetfilterHook {
	return &NetfilterHook{
		Hook:        HookInput,
		Tables:      []string{"mangle", "filter"},
		Description: "Packets destined for local delivery. Firewall rules (iptables -A INPUT) evaluated here.",
		Priority:    0,
	}
}

// NewForwardHook creates a netfilter FORWARD hook annotation.
// FORWARD is called for packets being routed through the machine.
func NewForwardHook() *NetfilterHook {
	return &NetfilterHook{
		Hook:        HookForward,
		Tables:      []string{"mangle", "filter"},
		Description: "Packets being forwarded/routed. Firewall rules (iptables -A FORWARD) evaluated here.",
		Priority:    0,
	}
}
