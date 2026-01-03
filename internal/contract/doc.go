// Package contract defines the data structures that model the Linux kernel
// networking stack for visualization purposes. This serves as the single
// source of truth for the frontend SPA.
//
// The contract models:
//   - sk_buff structure and its pointer arithmetic
//   - Kernel layer abstraction (User Space through Driver)
//   - Function call graph for egress/ingress paths
//   - State mutations as packets traverse the stack
//
// Based on Linux Kernel 5.10.8 as documented in:
// Stephan, A., & Wuestrich, L. (2024). The path of a packet through the linux kernel.
// Technical University of Munich, Chair of Network Architectures and Services.
package contract
