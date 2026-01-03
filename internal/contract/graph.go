package contract

// FunctionEdge represents a directed edge in the function call graph.
// It connects two functions and optionally includes a condition that
// determines when this path is taken.
type FunctionEdge struct {
	// From is the ID of the calling function
	From string `json:"from"`

	// To is the ID of the called function
	To string `json:"to"`

	// Condition describes when this edge is taken (empty for unconditional)
	// Examples: "TCP connection established", "No cached route", "Queue not full"
	Condition string `json:"condition,omitempty"`

	// IsErrorPath indicates if this edge represents an error handling path
	IsErrorPath bool `json:"isErrorPath,omitempty"`

	// Order is the sequence number for edges from the same source
	// Used to maintain consistent ordering in visualization
	Order int `json:"order,omitempty"`
}

// PacketPath represents a complete path through the kernel networking stack.
// This is the primary export format for the frontend visualization.
type PacketPath struct {
	// ID is a unique identifier for this path (e.g., "tcp_ipv4_egress")
	ID string `json:"id"`

	// Name is the display name for this path
	Name string `json:"name"`

	// Description explains what this path represents
	Description string `json:"description"`

	// Direction is either "egress" (sending) or "ingress" (receiving)
	Direction string `json:"direction"`

	// Protocol is the primary protocol of this path (e.g., "TCP", "UDP")
	Protocol string `json:"protocol"`

	// Functions is the list of all functions in this path
	Functions []KernelFunction `json:"functions"`

	// Edges defines the call relationships between functions
	Edges []FunctionEdge `json:"edges"`

	// EntryPoint is the ID of the starting function
	EntryPoint string `json:"entryPoint"`

	// ExitPoints are the IDs of possible ending functions
	ExitPoints []string `json:"exitPoints"`
}

// FunctionGraph is a helper structure for traversing the call graph.
type FunctionGraph struct {
	// functions maps function ID to function definition
	functions map[string]*KernelFunction

	// adjacency maps function ID to outgoing edges
	adjacency map[string][]FunctionEdge
}

// NewFunctionGraph creates a traversable graph from a PacketPath.
func NewFunctionGraph(path *PacketPath) *FunctionGraph {
	g := &FunctionGraph{
		functions: make(map[string]*KernelFunction),
		adjacency: make(map[string][]FunctionEdge),
	}

	for i := range path.Functions {
		f := &path.Functions[i]
		g.functions[f.ID] = f
	}

	for _, edge := range path.Edges {
		g.adjacency[edge.From] = append(g.adjacency[edge.From], edge)
	}

	return g
}

// GetFunction returns a function by ID.
func (g *FunctionGraph) GetFunction(id string) *KernelFunction {
	return g.functions[id]
}

// GetOutgoingEdges returns all edges originating from a function.
func (g *FunctionGraph) GetOutgoingEdges(id string) []FunctionEdge {
	return g.adjacency[id]
}

// GetNextFunctions returns the IDs of functions called by the given function.
func (g *FunctionGraph) GetNextFunctions(id string) []string {
	edges := g.adjacency[id]
	result := make([]string, len(edges))
	for i, edge := range edges {
		result[i] = edge.To
	}
	return result
}

// SimulateStep represents a single step in the packet simulation.
type SimulateStep struct {
	// StepNumber is the 1-indexed step number
	StepNumber int `json:"stepNumber"`

	// Function is the function being executed
	Function KernelFunction `json:"function"`

	// SKBuffState is the state of sk_buff after this step
	SKBuffState SKBuff `json:"skbuffState"`

	// EdgeTaken is the edge that led to this step (nil for entry point)
	EdgeTaken *FunctionEdge `json:"edgeTaken,omitempty"`

	// ConntrackState is the current connection tracking state (for TCP)
	ConntrackState *ConntrackEntry `json:"conntrackState,omitempty"`
}

// Simulate walks through the packet path and returns the sequence of steps.
// This is the core function that the frontend uses for animation.
func (path *PacketPath) Simulate(initialBufferSize int, payloadSize int) []SimulateStep {
	graph := NewFunctionGraph(path)
	steps := []SimulateStep{}

	// Initialize sk_buff with payload
	skb := NewSKBuffWithPayload(initialBufferSize, payloadSize)

	// Start at entry point
	currentID := path.EntryPoint
	stepNum := 1

	visited := make(map[string]bool)

	// For TCP data transfer, connection is already established
	conntrackState := NewConntrackEntry(ConntrackEstablished)

	for currentID != "" && !visited[currentID] {
		visited[currentID] = true

		fn := graph.GetFunction(currentID)
		if fn == nil {
			break
		}

		// Apply mutation if present
		if fn.SKBMutation != nil {
			switch fn.SKBMutation.Operation {
			case "push":
				skb.Push(fn.SKBMutation.HeaderType, fn.SKBMutation.Size)
			case "pull":
				skb.Pull(fn.SKBMutation.Size)
			case "put":
				skb.Put(fn.SKBMutation.Size)
			}
		}

		step := SimulateStep{
			StepNumber:     stepNum,
			Function:       *fn,
			SKBuffState:    *skb.Clone(),
			ConntrackState: conntrackState,
		}
		steps = append(steps, step)
		stepNum++

		// Get next function (take first non-error path for linear simulation)
		edges := graph.GetOutgoingEdges(currentID)
		currentID = ""
		for _, edge := range edges {
			if !edge.IsErrorPath {
				currentID = edge.To
				break
			}
		}
	}

	return steps
}

// SimulateIngress walks through the ingress path, starting with a full packet.
// Headers are progressively stripped (pulled) as the packet moves up the stack.
func (path *PacketPath) SimulateIngress(initialBufferSize int, payloadSize int) []SimulateStep {
	graph := NewFunctionGraph(path)
	steps := []SimulateStep{}

	// Initialize sk_buff with complete packet (all headers present)
	skb := NewSKBuffForIngress(initialBufferSize, payloadSize)

	// Start at entry point
	currentID := path.EntryPoint
	stepNum := 1

	visited := make(map[string]bool)

	// For TCP data reception, connection is already established
	conntrackState := NewConntrackEntry(ConntrackEstablished)

	for currentID != "" && !visited[currentID] {
		visited[currentID] = true

		fn := graph.GetFunction(currentID)
		if fn == nil {
			break
		}

		// Apply mutation if present
		if fn.SKBMutation != nil {
			switch fn.SKBMutation.Operation {
			case "push":
				skb.Push(fn.SKBMutation.HeaderType, fn.SKBMutation.Size)
			case "pull":
				skb.Pull(fn.SKBMutation.Size)
			case "put":
				skb.Put(fn.SKBMutation.Size)
			}
		}

		step := SimulateStep{
			StepNumber:     stepNum,
			Function:       *fn,
			SKBuffState:    *skb.Clone(),
			ConntrackState: conntrackState,
		}
		steps = append(steps, step)
		stepNum++

		// Get next function (take first non-error path for linear simulation)
		edges := graph.GetOutgoingEdges(currentID)
		currentID = ""
		for _, edge := range edges {
			if !edge.IsErrorPath {
				currentID = edge.To
				break
			}
		}
	}

	return steps
}
