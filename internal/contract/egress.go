package contract

// BuildTCPIPv4EgressPath constructs the complete TCP over IPv4 egress path
// based on Linux Kernel 5.10.8.
//
// This path represents a typical socket send operation using TCP,
// from the initial tcp_sendmsg call down to the NIC driver.
func BuildTCPIPv4EgressPath() *PacketPath {
	path := &PacketPath{
		ID:          "tcp_ipv4_egress",
		Name:        "TCP/IPv4 Egress Path",
		Description: "The path of a TCP packet from user space through the kernel to the network interface (Linux 5.10.8)",
		Direction:   "egress",
		Protocol:    "TCP",
		EntryPoint:  "tcp_sendmsg",
		ExitPoints:  []string{"ndo_start_xmit"},
	}

	// Define all functions in the egress path
	path.Functions = []KernelFunction{
		// Transport Layer - TCP
		{
			ID:           "tcp_sendmsg",
			Name:         "tcp_sendmsg",
			Layer:        LayerTransport,
			SourceFile:   "net/ipv4/tcp.c",
			LineNumber:   1434,
			Description:  "Entry point for TCP send operations. Acquires socket lock and delegates to tcp_sendmsg_locked.",
			IsEntryPoint: true,
		},
		{
			ID:          "tcp_sendmsg_locked",
			Name:        "tcp_sendmsg_locked",
			Layer:       LayerTransport,
			SourceFile:  "net/ipv4/tcp.c",
			LineNumber:  1172,
			Description: "Core TCP send logic. Allocates sk_buff and copies user data into kernel space.",
			SKBMutation: NewAllocMutation(2048, "Allocate sk_buff with headroom for all protocol headers"),
		},
		{
			ID:          "tcp_push",
			Name:        "tcp_push",
			Layer:       LayerTransport,
			SourceFile:  "net/ipv4/tcp.c",
			LineNumber:  689,
			Description: "Pushes pending data. Sets PSH flag if socket is being closed or buffer is full.",
		},
		{
			ID:          "__tcp_push_pending_frames",
			Name:        "__tcp_push_pending_frames",
			Layer:       LayerTransport,
			SourceFile:  "net/ipv4/tcp_output.c",
			LineNumber:  2556,
			Description: "Checks if there is data to send and initiates transmission.",
		},
		{
			ID:          "tcp_write_xmit",
			Name:        "tcp_write_xmit",
			Layer:       LayerTransport,
			SourceFile:  "net/ipv4/tcp_output.c",
			LineNumber:  2387,
			Description: "Main TCP transmission loop. Handles congestion control, pacing, and TSO segmentation.",
		},
		{
			ID:          "__tcp_transmit_skb",
			Name:        "__tcp_transmit_skb",
			Layer:       LayerTransport,
			SourceFile:  "net/ipv4/tcp_output.c",
			LineNumber:  1164,
			Description: "Builds the TCP header. Calculates checksum and sets sequence numbers.",
			SKBMutation: NewPushMutation("tcp", TCPHeaderSize),
		},

		// Network Layer - IP
		{
			ID:          "ip_queue_xmit",
			Name:        "ip_queue_xmit",
			Layer:       LayerNetwork,
			SourceFile:  "net/ipv4/ip_output.c",
			LineNumber:  470,
			Description: "Main IPv4 transmission entry point from transport layer. Handles routing lookup and IP header construction.",
			SKBMutation: NewPushMutation("ip", IPv4HeaderSize),
		},
		{
			ID:          "ip_local_out",
			Name:        "ip_local_out",
			Layer:       LayerNetwork,
			SourceFile:  "net/ipv4/ip_output.c",
			LineNumber:  115,
			Description: "Wrapper for locally generated packets. Calls __ip_local_out.",
		},
		{
			ID:            "__ip_local_out",
			Name:          "__ip_local_out",
			Layer:         LayerNetwork,
			SourceFile:    "net/ipv4/ip_output.c",
			LineNumber:    96,
			Description:   "Sets IP packet length and checksum. Invokes LOCAL_OUT netfilter hook.",
			NetfilterHook: NewOutputHook(),
		},
		{
			ID:            "ip_output",
			Name:          "ip_output",
			Layer:         LayerNetwork,
			SourceFile:    "net/ipv4/ip_output.c",
			LineNumber:    413,
			Description:   "Called after LOCAL_OUT hook. Invokes POST_ROUTING netfilter hook.",
			NetfilterHook: NewPostroutingHook(),
		},
		{
			ID:          "ip_finish_output",
			Name:        "ip_finish_output",
			Layer:       LayerNetwork,
			SourceFile:  "net/ipv4/ip_output.c",
			LineNumber:  311,
			Description: "BPF cgroup egress hook point. Handles GSO segmentation if needed.",
			BPFHook:     NewCgroupSKBHook("egress"),
		},
		{
			ID:          "__ip_finish_output",
			Name:        "__ip_finish_output",
			Layer:       LayerNetwork,
			SourceFile:  "net/ipv4/ip_output.c",
			LineNumber:  287,
			Description: "Checks MTU and fragments packet if necessary.",
		},
		{
			ID:          "ip_finish_output2",
			Name:        "ip_finish_output2",
			Layer:       LayerNetwork,
			SourceFile:  "net/ipv4/ip_output.c",
			LineNumber:  187,
			Description: "Resolves next-hop neighbor (ARP lookup) and prepares for L2 transmission.",
		},
		{
			ID:          "neigh_output",
			Name:        "neigh_output",
			Layer:       LayerNetwork,
			SourceFile:  "include/net/neighbour.h",
			LineNumber:  510,
			Description: "Neighbour subsystem output. Uses cached hardware header if available.",
		},
		{
			ID:          "neigh_hh_output",
			Name:        "neigh_hh_output",
			Layer:       LayerDataLink,
			SourceFile:  "include/net/neighbour.h",
			LineNumber:  490,
			Description: "Fast path using cached hardware header. Pushes Ethernet header.",
			SKBMutation: NewPushMutation("ethernet", EthernetHeaderSize),
		},

		// Data Link Layer - Queueing Discipline
		{
			ID:          "dev_queue_xmit",
			Name:        "dev_queue_xmit",
			Layer:       LayerDataLink,
			SourceFile:  "net/core/dev.c",
			LineNumber:  4044,
			Description: "Main device transmission entry point. Handles per-CPU processing.",
		},
		{
			ID:          "__dev_queue_xmit",
			Name:        "__dev_queue_xmit",
			Layer:       LayerDataLink,
			SourceFile:  "net/core/dev.c",
			LineNumber:  3954,
			Description: "Core queuing logic. TC egress BPF programs run here before qdisc.",
			BPFHook:     NewTCEgressHook(),
		},
		{
			ID:          "__dev_xmit_skb",
			Name:        "__dev_xmit_skb",
			Layer:       LayerDataLink,
			SourceFile:  "net/core/dev.c",
			LineNumber:  3683,
			Description: "Submits packet to qdisc. May queue or directly transmit based on qdisc state.",
		},
		{
			ID:          "sch_direct_xmit",
			Name:        "sch_direct_xmit",
			Layer:       LayerDataLink,
			SourceFile:  "net/sched/sch_generic.c",
			LineNumber:  310,
			Description: "Bypasses qdisc queue for direct transmission when possible.",
		},

		// Driver Layer
		{
			ID:          "dev_hard_start_xmit",
			Name:        "dev_hard_start_xmit",
			Layer:       LayerDriver,
			SourceFile:  "net/core/dev.c",
			LineNumber:  3506,
			Description: "Final generic layer before driver. Handles XDP and calls driver's ndo_start_xmit.",
		},
		{
			ID:          "ndo_start_xmit",
			Name:        "ndo_start_xmit",
			Layer:       LayerDriver,
			SourceFile:  "include/linux/netdevice.h",
			LineNumber:  1298,
			Description: "Driver-specific transmit function. Pointer to actual driver implementation (e.g., e1000, virtio-net).",
			IsExitPoint: true,
		},
	}

	// Define the edges (function call relationships)
	path.Edges = []FunctionEdge{
		{From: "tcp_sendmsg", To: "tcp_sendmsg_locked", Order: 1},
		{From: "tcp_sendmsg_locked", To: "tcp_push", Order: 1},
		{From: "tcp_push", To: "__tcp_push_pending_frames", Order: 1},
		{From: "__tcp_push_pending_frames", To: "tcp_write_xmit", Order: 1},
		{From: "tcp_write_xmit", To: "__tcp_transmit_skb", Order: 1},
		{From: "__tcp_transmit_skb", To: "ip_queue_xmit", Order: 1},
		{From: "ip_queue_xmit", To: "ip_local_out", Order: 1},
		{From: "ip_local_out", To: "__ip_local_out", Order: 1},
		{From: "__ip_local_out", To: "ip_output", Order: 1},
		{From: "ip_output", To: "ip_finish_output", Order: 1},
		{From: "ip_finish_output", To: "__ip_finish_output", Order: 1},
		{From: "__ip_finish_output", To: "ip_finish_output2", Order: 1},
		{From: "ip_finish_output2", To: "neigh_output", Order: 1},
		{From: "neigh_output", To: "neigh_hh_output", Order: 1, Condition: "Hardware header cached"},
		{From: "neigh_hh_output", To: "dev_queue_xmit", Order: 1},
		{From: "dev_queue_xmit", To: "__dev_queue_xmit", Order: 1},
		{From: "__dev_queue_xmit", To: "__dev_xmit_skb", Order: 1},
		{From: "__dev_xmit_skb", To: "sch_direct_xmit", Order: 1, Condition: "Direct transmit allowed"},
		{From: "sch_direct_xmit", To: "dev_hard_start_xmit", Order: 1},
		{From: "dev_hard_start_xmit", To: "ndo_start_xmit", Order: 1},
	}

	return path
}

// GetDefaultBufferSize returns the typical sk_buff allocation size
// that provides adequate headroom for all protocol headers.
func GetDefaultBufferSize() int {
	// Typical allocation: MTU (1500) + max headers + alignment
	// Ethernet: 14, IP: 60 (with options), TCP: 60 (with options)
	// Plus additional headroom for tunneling, etc.
	return 2048
}

// GetDefaultPayloadSize returns a typical payload size for simulation.
func GetDefaultPayloadSize() int {
	return 1000 // 1KB payload
}
