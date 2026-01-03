package contract

// ConntrackState represents a connection tracking state.
// Linux conntrack maintains state for stateful firewalling and NAT.
type ConntrackState string

// Connection tracking states for TCP
const (
	// ConntrackNew - New connection, no reply seen yet
	ConntrackNew ConntrackState = "NEW"

	// ConntrackSynSent - SYN packet sent, awaiting SYN-ACK
	ConntrackSynSent ConntrackState = "SYN_SENT"

	// ConntrackSynRecv - SYN received, SYN-ACK sent (server side)
	ConntrackSynRecv ConntrackState = "SYN_RECV"

	// ConntrackEstablished - Connection fully established
	ConntrackEstablished ConntrackState = "ESTABLISHED"

	// ConntrackFinWait - FIN sent, waiting for acknowledgment
	ConntrackFinWait ConntrackState = "FIN_WAIT"

	// ConntrackCloseWait - FIN received, waiting for close
	ConntrackCloseWait ConntrackState = "CLOSE_WAIT"

	// ConntrackLastAck - Final ACK expected
	ConntrackLastAck ConntrackState = "LAST_ACK"

	// ConntrackTimeWait - Waiting for stale packets to expire
	ConntrackTimeWait ConntrackState = "TIME_WAIT"

	// ConntrackClosed - Connection terminated
	ConntrackClosed ConntrackState = "CLOSED"
)

// ConntrackEntry represents the current connection tracking state
type ConntrackEntry struct {
	// State is the current conntrack state
	State ConntrackState `json:"state"`

	// Description explains the current state
	Description string `json:"description"`

	// Timeout is the remaining time before state expires (in seconds)
	Timeout int `json:"timeout,omitempty"`
}

// ConntrackStateDescriptions provides human-readable descriptions
var ConntrackStateDescriptions = map[ConntrackState]string{
	ConntrackNew:         "New connection. First packet seen, no reply yet.",
	ConntrackSynSent:     "SYN packet sent. Waiting for SYN-ACK from remote.",
	ConntrackSynRecv:     "SYN received, SYN-ACK sent. Awaiting final ACK.",
	ConntrackEstablished: "Connection established. Bidirectional traffic allowed.",
	ConntrackFinWait:     "FIN sent. Waiting for remote to acknowledge close.",
	ConntrackCloseWait:   "FIN received. Waiting for application to close.",
	ConntrackLastAck:     "Sent final FIN. Waiting for last ACK.",
	ConntrackTimeWait:    "Connection closed. Waiting for stale packets (2MSL).",
	ConntrackClosed:      "Connection fully closed. Entry will be removed.",
}

// NewConntrackEntry creates a conntrack entry with description
func NewConntrackEntry(state ConntrackState) *ConntrackEntry {
	return &ConntrackEntry{
		State:       state,
		Description: ConntrackStateDescriptions[state],
	}
}
