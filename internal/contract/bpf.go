package contract

// BPFHook represents an eBPF/XDP attachment point where BPF programs
// can intercept and modify packets.
type BPFHook struct {
	// Type is the BPF hook type: XDP, TC_INGRESS, TC_EGRESS, CGROUP_SKB, SOCKET
	Type string `json:"type"`

	// AttachPoint describes where the hook attaches in the kernel
	AttachPoint string `json:"attachPoint"`

	// Description explains what BPF programs can do at this hook
	Description string `json:"description"`

	// Actions lists the possible return values for this hook type
	Actions []string `json:"actions"`
}

// BPF hook type constants
const (
	BPFHookXDP       = "XDP"
	BPFHookTCIngress = "TC_INGRESS"
	BPFHookTCEgress  = "TC_EGRESS"
	BPFHookCgroupSKB = "CGROUP_SKB"
	BPFHookSocket    = "SOCKET"
)

// NewXDPHook creates an XDP hook annotation.
// XDP runs at the earliest point, before sk_buff allocation.
func NewXDPHook() *BPFHook {
	return &BPFHook{
		Type:        BPFHookXDP,
		AttachPoint: "NIC driver RX path",
		Description: "eXpress Data Path. Runs before sk_buff allocation for maximum performance. Can drop, pass, or redirect packets.",
		Actions:     []string{"XDP_PASS", "XDP_DROP", "XDP_TX", "XDP_REDIRECT", "XDP_ABORTED"},
	}
}

// NewTCIngressHook creates a TC ingress hook annotation.
// TC ingress runs after sk_buff is created, before the packet enters the stack.
func NewTCIngressHook() *BPFHook {
	return &BPFHook{
		Type:        BPFHookTCIngress,
		AttachPoint: "Traffic Control ingress qdisc",
		Description: "Traffic Control classifier. Can filter, modify, or redirect packets on ingress.",
		Actions:     []string{"TC_ACT_OK", "TC_ACT_SHOT", "TC_ACT_REDIRECT", "TC_ACT_PIPE"},
	}
}

// NewTCEgressHook creates a TC egress hook annotation.
// TC egress runs before the packet enters the qdisc for transmission.
func NewTCEgressHook() *BPFHook {
	return &BPFHook{
		Type:        BPFHookTCEgress,
		AttachPoint: "Traffic Control egress qdisc",
		Description: "Traffic Control classifier on egress. Can shape, filter, or redirect outgoing packets.",
		Actions:     []string{"TC_ACT_OK", "TC_ACT_SHOT", "TC_ACT_REDIRECT", "TC_ACT_PIPE"},
	}
}

// NewCgroupSKBHook creates a cgroup/skb hook annotation.
// Cgroup BPF is used for container networking policies.
func NewCgroupSKBHook(direction string) *BPFHook {
	return &BPFHook{
		Type:        BPFHookCgroupSKB,
		AttachPoint: "Cgroup " + direction + " path",
		Description: "Cgroup socket buffer hook. Used for container networking policies and egress filtering.",
		Actions:     []string{"ALLOW", "DENY"},
	}
}

// NewSocketBPFHook creates a socket-level BPF hook annotation.
func NewSocketBPFHook() *BPFHook {
	return &BPFHook{
		Type:        BPFHookSocket,
		AttachPoint: "Socket layer",
		Description: "Socket-level BPF. Can filter packets before they reach the application.",
		Actions:     []string{"ALLOW", "DENY"},
	}
}
