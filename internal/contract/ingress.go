package contract

// BuildTCPIPv4IngressPath constructs the complete TCP over IPv4 ingress path
// based on Linux Kernel 5.10.8.
//
// This path represents a typical packet reception from the NIC driver
// up through NAPI, the network stack, to the socket layer.
func BuildTCPIPv4IngressPath() *PacketPath {
	path := &PacketPath{
		ID:          "tcp_ipv4_ingress",
		Name:        "TCP/IPv4 Ingress Path",
		Description: "The path of a TCP packet from the network interface through the kernel to user space (Linux 5.10.8)",
		Direction:   "ingress",
		Protocol:    "TCP",
		EntryPoint:  "napi_poll",
		ExitPoints:  []string{"sk_data_ready"},
	}

	// Define all functions in the ingress path
	path.Functions = []KernelFunction{
		// Driver Layer - NAPI
		{
			ID:           "napi_poll",
			Name:         "napi_poll",
			Layer:        LayerDriver,
			SourceFile:   "net/core/dev.c",
			LineNumber:   6740,
			Description:  "NAPI polling entry point. Called by softirq to process received packets from the driver's ring buffer.",
			IsEntryPoint: true,
		},
		{
			ID:          "napi_gro_receive",
			Name:        "napi_gro_receive",
			Layer:       LayerDriver,
			SourceFile:  "net/core/dev.c",
			LineNumber:  6081,
			Description: "Generic Receive Offload handler. XDP programs run here before sk_buff allocation.",
			BPFHook:     NewXDPHook(),
		},
		{
			ID:          "napi_skb_finish",
			Name:        "napi_skb_finish",
			Layer:       LayerDriver,
			SourceFile:  "net/core/dev.c",
			LineNumber:  6052,
			Description: "Finishes GRO processing and passes the sk_buff up the stack.",
		},

		// Data Link Layer
		{
			ID:          "netif_receive_skb",
			Name:        "netif_receive_skb",
			Layer:       LayerDataLink,
			SourceFile:  "net/core/dev.c",
			LineNumber:  5583,
			Description: "Main entry point for receiving packets from the driver. Timestamps and prepares the packet.",
		},
		{
			ID:          "netif_receive_skb_internal",
			Name:        "netif_receive_skb_internal",
			Layer:       LayerDataLink,
			SourceFile:  "net/core/dev.c",
			LineNumber:  5508,
			Description: "Internal receive handler. Handles RPS (Receive Packet Steering) if enabled.",
		},
		{
			ID:          "__netif_receive_skb",
			Name:        "__netif_receive_skb",
			Layer:       LayerDataLink,
			SourceFile:  "net/core/dev.c",
			LineNumber:  5405,
			Description: "Core receive function. TC ingress BPF programs and generic XDP run here.",
			BPFHook:     NewTCIngressHook(),
		},
		{
			ID:          "__netif_receive_skb_one_core",
			Name:        "__netif_receive_skb_one_core",
			Layer:       LayerDataLink,
			SourceFile:  "net/core/dev.c",
			LineNumber:  5303,
			Description: "Single-core receive path. Processes packet on current CPU.",
		},
		{
			ID:          "__netif_receive_skb_core",
			Name:        "__netif_receive_skb_core",
			Layer:       LayerDataLink,
			SourceFile:  "net/core/dev.c",
			LineNumber:  5099,
			Description: "Core packet classification. Strips Ethernet header and determines protocol handler.",
			SKBMutation: NewPullMutation("ethernet", EthernetHeaderSize),
		},
		{
			ID:          "deliver_skb",
			Name:        "deliver_skb",
			Layer:       LayerDataLink,
			SourceFile:  "net/core/dev.c",
			LineNumber:  2248,
			Description: "Delivers packet to the registered protocol handler (e.g., ip_rcv for IPv4).",
		},

		// Network Layer - IP
		{
			ID:            "ip_rcv",
			Name:          "ip_rcv",
			Layer:         LayerNetwork,
			SourceFile:    "net/ipv4/ip_input.c",
			LineNumber:    530,
			Description:   "IPv4 receive entry point. Validates IP header checksum and invokes PREROUTING netfilter hook.",
			NetfilterHook: NewPreroutingHook(),
		},
		{
			ID:          "ip_rcv_finish",
			Name:        "ip_rcv_finish",
			Layer:       LayerNetwork,
			SourceFile:  "net/ipv4/ip_input.c",
			LineNumber:  414,
			Description: "Finishes IP header processing. Performs routing lookup and strips IP header.",
			SKBMutation: NewPullMutation("ip", IPv4HeaderSize),
		},
		{
			ID:          "ip_local_deliver",
			Name:        "ip_local_deliver",
			Layer:       LayerNetwork,
			SourceFile:  "net/ipv4/ip_input.c",
			LineNumber:  240,
			Description: "Handles locally destined packets. Reassembles IP fragments if needed.",
		},
		{
			ID:            "ip_local_deliver_finish",
			Name:          "ip_local_deliver_finish",
			Layer:         LayerNetwork,
			SourceFile:    "net/ipv4/ip_input.c",
			LineNumber:    226,
			Description:   "Invokes INPUT netfilter hook before passing to transport layer.",
			NetfilterHook: NewInputHook(),
		},
		{
			ID:          "ip_protocol_deliver_rcu",
			Name:        "ip_protocol_deliver_rcu",
			Layer:       LayerNetwork,
			SourceFile:  "net/ipv4/ip_input.c",
			LineNumber:  187,
			Description: "Dispatches packet to the transport protocol handler based on IP protocol field.",
		},

		// Transport Layer - TCP
		{
			ID:          "tcp_v4_rcv",
			Name:        "tcp_v4_rcv",
			Layer:       LayerTransport,
			SourceFile:  "net/ipv4/tcp_ipv4.c",
			LineNumber:  1915,
			Description: "TCP receive entry point. Validates TCP checksum and looks up socket.",
		},
		{
			ID:          "tcp_v4_do_rcv",
			Name:        "tcp_v4_do_rcv",
			Layer:       LayerTransport,
			SourceFile:  "net/ipv4/tcp_ipv4.c",
			LineNumber:  1655,
			Description: "Main TCP receive handler. Processes TCP header and updates connection state.",
			SKBMutation: NewPullMutation("tcp", TCPHeaderSize),
		},
		{
			ID:          "tcp_rcv_established",
			Name:        "tcp_rcv_established",
			Layer:       LayerTransport,
			SourceFile:  "net/ipv4/tcp_input.c",
			LineNumber:  5704,
			Description: "Fast path for established connections. Handles ACKs, window updates, and data.",
		},
		{
			ID:          "tcp_data_queue",
			Name:        "tcp_data_queue",
			Layer:       LayerTransport,
			SourceFile:  "net/ipv4/tcp_input.c",
			LineNumber:  4919,
			Description: "Queues received data. Handles out-of-order segments and SACK.",
		},
		{
			ID:          "tcp_queue_rcv",
			Name:        "tcp_queue_rcv",
			Layer:       LayerTransport,
			SourceFile:  "net/ipv4/tcp_input.c",
			LineNumber:  4837,
			Description: "Adds data to socket receive queue. Updates TCP receive window.",
		},

		// Socket Layer
		{
			ID:          "sk_data_ready",
			Name:        "sk_data_ready",
			Layer:       LayerSocket,
			SourceFile:  "net/core/sock.c",
			LineNumber:  2990,
			Description: "Wakes up any process waiting to read from the socket. Data is now available for recv().",
			IsExitPoint: true,
		},
	}

	// Define the edges (function call relationships)
	path.Edges = []FunctionEdge{
		{From: "napi_poll", To: "napi_gro_receive", Order: 1},
		{From: "napi_gro_receive", To: "napi_skb_finish", Order: 1},
		{From: "napi_skb_finish", To: "netif_receive_skb", Order: 1},
		{From: "netif_receive_skb", To: "netif_receive_skb_internal", Order: 1},
		{From: "netif_receive_skb_internal", To: "__netif_receive_skb", Order: 1},
		{From: "__netif_receive_skb", To: "__netif_receive_skb_one_core", Order: 1},
		{From: "__netif_receive_skb_one_core", To: "__netif_receive_skb_core", Order: 1},
		{From: "__netif_receive_skb_core", To: "deliver_skb", Order: 1},
		{From: "deliver_skb", To: "ip_rcv", Order: 1, Condition: "Protocol is IPv4"},
		{From: "ip_rcv", To: "ip_rcv_finish", Order: 1},
		{From: "ip_rcv_finish", To: "ip_local_deliver", Order: 1, Condition: "Destination is local"},
		{From: "ip_local_deliver", To: "ip_local_deliver_finish", Order: 1},
		{From: "ip_local_deliver_finish", To: "ip_protocol_deliver_rcu", Order: 1},
		{From: "ip_protocol_deliver_rcu", To: "tcp_v4_rcv", Order: 1, Condition: "Protocol is TCP"},
		{From: "tcp_v4_rcv", To: "tcp_v4_do_rcv", Order: 1, Condition: "Socket found"},
		{From: "tcp_v4_do_rcv", To: "tcp_rcv_established", Order: 1, Condition: "Connection established"},
		{From: "tcp_rcv_established", To: "tcp_data_queue", Order: 1, Condition: "Has data"},
		{From: "tcp_data_queue", To: "tcp_queue_rcv", Order: 1},
		{From: "tcp_queue_rcv", To: "sk_data_ready", Order: 1},
	}

	return path
}

// NewSKBuffForIngress creates an sk_buff as it would appear when received from the NIC.
// The buffer contains the full packet with all headers already present.
func NewSKBuffForIngress(totalSize, payloadSize int) *SKBuff {
	// For ingress, the packet arrives complete with all headers
	// Data starts at 0 (beginning of buffer) with all headers present
	headerSize := EthernetHeaderSize + IPv4HeaderSize + TCPHeaderSize
	totalPacketLen := headerSize + payloadSize

	skb := &SKBuff{
		Head: 0,
		Data: 0,
		Tail: totalPacketLen,
		End:  totalSize,
		Layers: []ProtocolHeader{
			{Protocol: "ethernet", Offset: 0, Size: EthernetHeaderSize},
			{Protocol: "ip", Offset: EthernetHeaderSize, Size: IPv4HeaderSize},
			{Protocol: "tcp", Offset: EthernetHeaderSize + IPv4HeaderSize, Size: TCPHeaderSize},
		},
	}

	return skb
}
