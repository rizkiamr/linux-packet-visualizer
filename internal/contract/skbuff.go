package contract

// SKBuff represents the Linux kernel's sk_buff structure, which is the
// fundamental data structure for network packet handling. It models the
// memory layout with four critical pointers that define the packet boundaries.
//
// Memory Layout:
//
//	+------------------+
//	|    headroom      |  <- space for prepending headers
//	+------------------+ <- Head
//	|                  |
//	|   packet data    |  <- Data to Tail
//	|                  |
//	+------------------+ <- Tail
//	|    tailroom      |  <- space for appending data
//	+------------------+ <- End
//
// During egress (sending), headers are pushed onto the front of the packet,
// moving the Data pointer backward toward Head.
//
// During ingress (receiving), headers are pulled from the front,
// moving the Data pointer forward toward Tail.
type SKBuff struct {
	// Head points to the start of the allocated buffer.
	// This value is immutable after allocation.
	Head int `json:"head"`

	// Data points to the start of the actual packet data.
	// This pointer moves backward when headers are pushed (egress)
	// and forward when headers are pulled (ingress).
	Data int `json:"data"`

	// Tail points to the end of the actual packet data.
	// This pointer moves forward when data is appended.
	Tail int `json:"tail"`

	// End points to the end of the allocated buffer.
	// This value is immutable after allocation.
	End int `json:"end"`

	// Layers tracks which protocol headers are currently present
	// in the buffer, in order from outermost to innermost.
	Layers []ProtocolHeader `json:"layers"`
}

// ProtocolHeader represents a single protocol header within the sk_buff.
type ProtocolHeader struct {
	// Protocol identifies the header type (e.g., "ethernet", "ip", "tcp")
	Protocol string `json:"protocol"`

	// Offset is the byte offset from the current Data pointer
	Offset int `json:"offset"`

	// Size is the header size in bytes
	Size int `json:"size"`
}

// NewSKBuff creates a new sk_buff with the specified total buffer size.
// Initially, Data and Tail point to the same location (empty payload),
// with maximum headroom for header prepending.
func NewSKBuff(totalSize int) *SKBuff {
	// Start with Data/Tail at the end to allow maximum header space
	// This is typical for egress path where we build headers backward
	return &SKBuff{
		Head:   0,
		Data:   totalSize,
		Tail:   totalSize,
		End:    totalSize,
		Layers: []ProtocolHeader{},
	}
}

// NewSKBuffWithPayload creates an sk_buff with an initial payload.
// The payload is placed at the end of the buffer, leaving headroom
// for protocol headers to be pushed during egress.
func NewSKBuffWithPayload(totalSize, payloadSize int) *SKBuff {
	dataStart := totalSize - payloadSize
	return &SKBuff{
		Head:   0,
		Data:   dataStart,
		Tail:   totalSize,
		End:    totalSize,
		Layers: []ProtocolHeader{},
	}
}

// Push prepends space for a header at the front of the packet.
// This moves the Data pointer backward by the specified size.
// Returns false if there is insufficient headroom.
func (s *SKBuff) Push(protocol string, size int) bool {
	newData := s.Data - size
	if newData < s.Head {
		return false // insufficient headroom
	}
	s.Data = newData

	// Add the header to the front of the layers list
	header := ProtocolHeader{
		Protocol: protocol,
		Offset:   0,
		Size:     size,
	}

	// Update offsets of existing headers
	for i := range s.Layers {
		s.Layers[i].Offset += size
	}

	s.Layers = append([]ProtocolHeader{header}, s.Layers...)
	return true
}

// Pull removes a header from the front of the packet.
// This moves the Data pointer forward by the specified size.
// Returns false if the pull would exceed the Tail pointer.
func (s *SKBuff) Pull(size int) bool {
	newData := s.Data + size
	if newData > s.Tail {
		return false // would exceed packet data
	}
	s.Data = newData

	// Remove the first header and update offsets
	if len(s.Layers) > 0 {
		s.Layers = s.Layers[1:]
		for i := range s.Layers {
			s.Layers[i].Offset -= size
		}
	}
	return true
}

// Put appends data to the end of the packet.
// This moves the Tail pointer forward by the specified size.
// Returns false if there is insufficient tailroom.
func (s *SKBuff) Put(size int) bool {
	newTail := s.Tail + size
	if newTail > s.End {
		return false // insufficient tailroom
	}
	s.Tail = newTail
	return true
}

// Headroom returns the available space before the Data pointer.
func (s *SKBuff) Headroom() int {
	return s.Data - s.Head
}

// Tailroom returns the available space after the Tail pointer.
func (s *SKBuff) Tailroom() int {
	return s.End - s.Tail
}

// Len returns the current packet length (Data to Tail).
func (s *SKBuff) Len() int {
	return s.Tail - s.Data
}

// Clone creates a deep copy of the sk_buff.
func (s *SKBuff) Clone() *SKBuff {
	clone := &SKBuff{
		Head:   s.Head,
		Data:   s.Data,
		Tail:   s.Tail,
		End:    s.End,
		Layers: make([]ProtocolHeader, len(s.Layers)),
	}
	copy(clone.Layers, s.Layers)
	return clone
}
