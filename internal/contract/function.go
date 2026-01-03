package contract

// KernelFunction represents a single function node in the kernel call graph.
// Each function has metadata about its location, purpose, and how it
// mutates the sk_buff structure.
type KernelFunction struct {
	// ID is a unique identifier for the function (e.g., "tcp_sendmsg")
	ID string `json:"id"`

	// Name is the display name shown in the visualization
	Name string `json:"name"`

	// Layer indicates which kernel layer this function belongs to
	Layer Layer `json:"layer"`

	// SourceFile is the kernel source file path (e.g., "net/ipv4/tcp.c")
	SourceFile string `json:"sourceFile"`

	// LineNumber is the approximate line number in the kernel source (5.10.8)
	LineNumber int `json:"lineNumber,omitempty"`

	// Description is a brief explanation of what the function does
	Description string `json:"description"`

	// SKBMutation describes how this function modifies the sk_buff (nil if no change)
	SKBMutation *SKBMutation `json:"skbMutation,omitempty"`

	// NetfilterHook indicates if this function triggers a netfilter hook (nil if none)
	NetfilterHook *NetfilterHook `json:"netfilterHook,omitempty"`

	// BPFHook indicates if this function has a BPF/XDP attachment point (nil if none)
	BPFHook *BPFHook `json:"bpfHook,omitempty"`

	// IsEntryPoint indicates if this is a valid starting point for a path
	IsEntryPoint bool `json:"isEntryPoint,omitempty"`

	// IsExitPoint indicates if this is an endpoint (packet leaves kernel)
	IsExitPoint bool `json:"isExitPoint,omitempty"`
}

// SKBMutation describes how a function modifies the sk_buff structure.
type SKBMutation struct {
	// Operation is the type of mutation: "push", "pull", "put", "alloc", "free"
	Operation string `json:"operation"`

	// HeaderType is the protocol header affected (e.g., "tcp", "ip", "ethernet")
	HeaderType string `json:"headerType,omitempty"`

	// Size is the number of bytes affected by the operation
	Size int `json:"size"`

	// Description is a human-readable explanation of the mutation
	Description string `json:"description"`
}

// Common header sizes in bytes
const (
	// EthernetHeaderSize is the standard Ethernet II header size (no VLAN)
	EthernetHeaderSize = 14

	// IPv4HeaderSize is the minimum IPv4 header size (no options)
	IPv4HeaderSize = 20

	// IPv6HeaderSize is the fixed IPv6 header size
	IPv6HeaderSize = 40

	// TCPHeaderSize is the minimum TCP header size (no options)
	TCPHeaderSize = 20

	// UDPHeaderSize is the fixed UDP header size
	UDPHeaderSize = 8

	// ICMPHeaderSize is the minimum ICMP header size
	ICMPHeaderSize = 8
)

// NewPushMutation creates a mutation representing a header push operation.
func NewPushMutation(headerType string, size int) *SKBMutation {
	return &SKBMutation{
		Operation:   "push",
		HeaderType:  headerType,
		Size:        size,
		Description: "Push " + headerType + " header",
	}
}

// NewPullMutation creates a mutation representing a header pull operation.
func NewPullMutation(headerType string, size int) *SKBMutation {
	return &SKBMutation{
		Operation:   "pull",
		HeaderType:  headerType,
		Size:        size,
		Description: "Pull " + headerType + " header",
	}
}

// NewAllocMutation creates a mutation representing sk_buff allocation.
func NewAllocMutation(size int, description string) *SKBMutation {
	return &SKBMutation{
		Operation:   "alloc",
		Size:        size,
		Description: description,
	}
}
