/**
 * FunctionNode - Renders a single kernel function as an SVG group.
 */
export function FunctionNode({
    id,
    name,
    layer,
    sourceFile,
    description,
    skbMutation,
    netfilterHook,
    bpfHook,
    x,
    y,
    width = 150,
    height = 50,
    isActive = false,
    isTraversed = false,
    onClick
}) {
    // Map layer to CSS class
    const layerClass = {
        'User Space': 'layer-user',
        'Socket Layer': 'layer-socket',
        'Transport Layer': 'layer-transport',
        'Network Layer': 'layer-network',
        'Data Link Layer': 'layer-datalink',
        'Device Driver': 'layer-driver',
    }[layer] || 'layer-transport';

    // Netfilter hook colors (warm colors)
    const nfHookColors = {
        'PREROUTING': '#ff6b6b',
        'INPUT': '#4ecdc4',
        'OUTPUT': '#ffe66d',
        'POSTROUTING': '#95e1d3',
        'FORWARD': '#a8e6cf',
    };

    // BPF hook colors (cool/purple colors)
    const bpfHookColors = {
        'XDP': '#ff00ff',
        'TC_INGRESS': '#00bfff',
        'TC_EGRESS': '#1e90ff',
        'CGROUP_SKB': '#9370db',
        'SOCKET': '#da70d6',
    };

    const hasHook = netfilterHook || bpfHook;

    const classes = [
        'function-node',
        layerClass,
        isActive && 'active',
        isTraversed && 'traversed',
        hasHook && 'has-hook',
        bpfHook && 'has-bpf',
    ].filter(Boolean).join(' ');

    // Extract just the filename from the source path
    const fileName = sourceFile?.split('/').pop() || '';

    // Calculate badge offset - if both hooks, show them stacked
    const hasBothHooks = netfilterHook && bpfHook;
    const nodeHeight = hasBothHooks ? height + 38 : (hasHook ? height + 20 : height);

    return (
        <g
            className={classes}
            transform={`translate(${x}, ${y})`}
            onClick={() => onClick?.(id)}
        >
            {/* Background rectangle */}
            <rect
                className="node-bg"
                x={0}
                y={0}
                width={width}
                height={nodeHeight}
                rx={8}
                ry={8}
            />

            {/* Function name */}
            <text
                className="node-name"
                x={width / 2}
                y={20}
                textAnchor="middle"
                dominantBaseline="middle"
            >
                {name}
            </text>

            {/* Source file */}
            <text
                className="node-file"
                x={width / 2}
                y={38}
                textAnchor="middle"
                dominantBaseline="middle"
            >
                {fileName}
            </text>

            {/* Mutation indicator */}
            {skbMutation && (
                <circle
                    className="mutation-indicator"
                    cx={width - 12}
                    cy={12}
                    r={5}
                >
                    <title>{skbMutation.description}</title>
                </circle>
            )}

            {/* Netfilter hook badge (rectangular) */}
            {netfilterHook && (
                <g className="netfilter-badge" transform={`translate(${width / 2}, ${height + 2})`}>
                    <rect
                        x={-35}
                        y={0}
                        width={70}
                        height={16}
                        rx={3}
                        fill={nfHookColors[netfilterHook.hook] || '#888'}
                        opacity={0.9}
                    />
                    <text
                        x={0}
                        y={11}
                        textAnchor="middle"
                        fill="#000"
                        fontSize={9}
                        fontWeight={600}
                        fontFamily="var(--font-mono)"
                    >
                        {netfilterHook.hook}
                    </text>
                </g>
            )}

            {/* BPF hook badge (pill-shaped) */}
            {bpfHook && (
                <g className="bpf-badge" transform={`translate(${width / 2}, ${netfilterHook ? height + 20 : height + 2})`}>
                    <rect
                        x={-40}
                        y={0}
                        width={80}
                        height={16}
                        rx={8}
                        fill={bpfHookColors[bpfHook.type] || '#888'}
                        opacity={0.95}
                    />
                    <text
                        x={0}
                        y={11}
                        textAnchor="middle"
                        fill="#fff"
                        fontSize={9}
                        fontWeight={700}
                        fontFamily="var(--font-mono)"
                    >
                        {bpfHook.type.replace('_', ' ')}
                    </text>
                </g>
            )}
        </g>
    );
}
