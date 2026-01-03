package contract

// Layer represents a layer in the Linux kernel networking stack.
// These correspond to the visual tiers in the frontend layout.
type Layer int

const (
	// LayerUserSpace represents the user space syscall interface.
	// Functions: write(), sendto(), sendmsg()
	LayerUserSpace Layer = iota

	// LayerSocket represents the socket abstraction layer.
	// This is where the socket API interfaces with protocol-specific code.
	LayerSocket

	// LayerTransport represents the transport layer (L4).
	// Protocols: TCP, UDP, SCTP
	// Functions: tcp_sendmsg, udp_sendmsg, etc.
	LayerTransport

	// LayerNetwork represents the network layer (L3).
	// Protocols: IPv4, IPv6
	// Functions: ip_queue_xmit, ip_local_out, ip_output, etc.
	LayerNetwork

	// LayerDataLink represents the data link layer (L2).
	// This includes the queueing discipline (qdisc) and neighbor subsystem.
	// Functions: dev_queue_xmit, neigh_output, etc.
	LayerDataLink

	// LayerDriver represents the network device driver layer.
	// This is where the packet is handed to the NIC hardware.
	// Functions: dev_hard_start_xmit, ndo_start_xmit
	LayerDriver
)

// String returns the human-readable name of the layer.
func (l Layer) String() string {
	switch l {
	case LayerUserSpace:
		return "User Space"
	case LayerSocket:
		return "Socket Layer"
	case LayerTransport:
		return "Transport Layer"
	case LayerNetwork:
		return "Network Layer"
	case LayerDataLink:
		return "Data Link Layer"
	case LayerDriver:
		return "Device Driver"
	default:
		return "Unknown"
	}
}

// CSSClass returns a CSS-friendly class name for the layer.
func (l Layer) CSSClass() string {
	switch l {
	case LayerUserSpace:
		return "layer-user"
	case LayerSocket:
		return "layer-socket"
	case LayerTransport:
		return "layer-transport"
	case LayerNetwork:
		return "layer-network"
	case LayerDataLink:
		return "layer-datalink"
	case LayerDriver:
		return "layer-driver"
	default:
		return "layer-unknown"
	}
}

// MarshalJSON implements custom JSON marshaling for Layer.
func (l Layer) MarshalJSON() ([]byte, error) {
	return []byte(`"` + l.String() + `"`), nil
}

// UnmarshalJSON implements custom JSON unmarshaling for Layer.
func (l *Layer) UnmarshalJSON(data []byte) error {
	// Remove quotes from the string
	s := string(data)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	switch s {
	case "User Space":
		*l = LayerUserSpace
	case "Socket Layer":
		*l = LayerSocket
	case "Transport Layer":
		*l = LayerTransport
	case "Network Layer":
		*l = LayerNetwork
	case "Data Link Layer":
		*l = LayerDataLink
	case "Device Driver":
		*l = LayerDriver
	default:
		*l = LayerUserSpace // Default fallback
	}
	return nil
}
