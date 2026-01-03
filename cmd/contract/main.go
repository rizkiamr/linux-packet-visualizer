// Command contract generates the JSON data contract for the Linux Packet Visualizer.
//
// Usage:
//
//	go run ./cmd/contract > egress_path.json
//	go run ./cmd/contract -o frontend/public/data/egress_path.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rzkiamr/linux-packet-visualizer/internal/contract"
)

func main() {
	outputFile := flag.String("o", "", "Output file path (default: stdout)")
	compact := flag.Bool("compact", false, "Output compact JSON (no indentation)")
	noSim := flag.Bool("no-sim", false, "Exclude pre-computed simulation")
	bufferSize := flag.Int("buffer", 2048, "sk_buff buffer size for simulation")
	payloadSize := flag.Int("payload", 1000, "Initial payload size for simulation")

	flag.Parse()

	opts := contract.ExportOptions{
		Pretty:            !*compact,
		IncludeSimulation: !*noSim,
		BufferSize:        *bufferSize,
		PayloadSize:       *payloadSize,
	}

	data, err := contract.ExportTCPIPv4EgressPath(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating contract: %v\n", err)
		os.Exit(1)
	}

	// Add timestamp to the output
	var export contract.ExportPacket
	if err := json.Unmarshal(data, &export); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing generated contract: %v\n", err)
		os.Exit(1)
	}
	export.GeneratedAt = time.Now().UTC().Format(time.RFC3339)

	if opts.Pretty {
		data, err = json.MarshalIndent(export, "", "  ")
	} else {
		data, err = json.Marshal(export)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error re-marshaling contract: %v\n", err)
		os.Exit(1)
	}

	if *outputFile != "" {
		if err := os.WriteFile(*outputFile, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Contract written to %s\n", *outputFile)
	} else {
		fmt.Println(string(data))
	}
}
