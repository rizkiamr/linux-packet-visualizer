package contract

import (
	"encoding/json"
)

// ExportOptions configures the JSON export.
type ExportOptions struct {
	// Pretty enables indented JSON output
	Pretty bool

	// IncludeSimulation includes a pre-computed simulation run
	IncludeSimulation bool

	// BufferSize is the sk_buff size for simulation (default: 2048)
	BufferSize int

	// PayloadSize is the initial payload size for simulation (default: 1000)
	PayloadSize int
}

// DefaultExportOptions returns sensible defaults for export.
func DefaultExportOptions() ExportOptions {
	return ExportOptions{
		Pretty:            true,
		IncludeSimulation: true,
		BufferSize:        GetDefaultBufferSize(),
		PayloadSize:       GetDefaultPayloadSize(),
	}
}

// ExportPacket is the complete export structure for frontend consumption.
// Supports multiple paths (egress and ingress).
type ExportPacket struct {
	// Version is the contract schema version
	Version string `json:"version"`

	// KernelVersion is the Linux kernel version this is based on
	KernelVersion string `json:"kernelVersion"`

	// GeneratedAt is the generation timestamp
	GeneratedAt string `json:"generatedAt"`

	// Paths contains all available packet paths
	Paths []PathWithSimulation `json:"paths"`

	// Metadata contains additional information for the frontend
	Metadata ExportMetadata `json:"metadata"`
}

// PathWithSimulation bundles a path with its pre-computed simulation.
type PathWithSimulation struct {
	// Path is the packet path definition
	Path PacketPath `json:"path"`

	// Simulation is the pre-computed simulation (optional)
	Simulation []SimulateStep `json:"simulation,omitempty"`
}

// ExportMetadata contains frontend-relevant metadata.
type ExportMetadata struct {
	// Layers lists all layers in order for rendering
	Layers []LayerInfo `json:"layers"`

	// HeaderSizes maps protocol names to their sizes
	HeaderSizes map[string]int `json:"headerSizes"`

	// BufferSize is the total sk_buff size used
	BufferSize int `json:"bufferSize"`

	// PayloadSize is the initial payload size
	PayloadSize int `json:"payloadSize"`
}

// LayerInfo provides rendering information for a layer.
type LayerInfo struct {
	// ID is the layer identifier
	ID string `json:"id"`

	// Name is the display name
	Name string `json:"name"`

	// CSSClass is the CSS class for styling
	CSSClass string `json:"cssClass"`

	// Order is the rendering order (0 = top)
	Order int `json:"order"`
}

// ExportAllPaths exports both egress and ingress paths as JSON.
func ExportAllPaths(opts ExportOptions) ([]byte, error) {
	egressPath := BuildTCPIPv4EgressPath()
	ingressPath := BuildTCPIPv4IngressPath()

	paths := []PathWithSimulation{
		{Path: *egressPath},
		{Path: *ingressPath},
	}

	if opts.IncludeSimulation {
		// Egress simulation: start with payload, push headers
		paths[0].Simulation = egressPath.Simulate(opts.BufferSize, opts.PayloadSize)

		// Ingress simulation: start with full packet, pull headers
		paths[1].Simulation = ingressPath.SimulateIngress(opts.BufferSize, opts.PayloadSize)
	}

	export := ExportPacket{
		Version:       "1.1.0",
		KernelVersion: "5.10.8",
		GeneratedAt:   "", // Will be set by caller if needed
		Paths:         paths,
		Metadata: ExportMetadata{
			Layers: []LayerInfo{
				{ID: "user", Name: "User Space", CSSClass: "layer-user", Order: 0},
				{ID: "socket", Name: "Socket Layer", CSSClass: "layer-socket", Order: 1},
				{ID: "transport", Name: "Transport Layer", CSSClass: "layer-transport", Order: 2},
				{ID: "network", Name: "Network Layer", CSSClass: "layer-network", Order: 3},
				{ID: "datalink", Name: "Data Link Layer", CSSClass: "layer-datalink", Order: 4},
				{ID: "driver", Name: "Device Driver", CSSClass: "layer-driver", Order: 5},
			},
			HeaderSizes: map[string]int{
				"ethernet": EthernetHeaderSize,
				"ip":       IPv4HeaderSize,
				"ipv6":     IPv6HeaderSize,
				"tcp":      TCPHeaderSize,
				"udp":      UDPHeaderSize,
				"icmp":     ICMPHeaderSize,
			},
			BufferSize:  opts.BufferSize,
			PayloadSize: opts.PayloadSize,
		},
	}

	if opts.Pretty {
		return json.MarshalIndent(export, "", "  ")
	}
	return json.Marshal(export)
}

// ExportAllPathsJSON is a convenience function with default options.
func ExportAllPathsJSON() ([]byte, error) {
	return ExportAllPaths(DefaultExportOptions())
}

// Legacy: ExportTCPIPv4EgressPath exports only the egress path (for backward compatibility).
func ExportTCPIPv4EgressPath(opts ExportOptions) ([]byte, error) {
	return ExportAllPaths(opts)
}

// Legacy: ExportEgressPathJSON is a convenience function with default options.
func ExportEgressPathJSON() ([]byte, error) {
	return ExportAllPathsJSON()
}
