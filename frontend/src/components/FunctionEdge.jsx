/**
 * FunctionEdge - Renders a directed edge between two function nodes.
 */
export function FunctionEdge({
    id,
    fromX,
    fromY,
    toX,
    toY,
    isActive = false,
    isTraversed = false
}) {
    // Calculate control points for a smooth curve
    // Edges typically go downward, so we create a nice vertical curve
    const midY = (fromY + toY) / 2;
    const path = `M ${fromX} ${fromY} C ${fromX} ${midY}, ${toX} ${midY}, ${toX} ${toY}`;

    const classes = [
        'function-edge',
        isActive && 'active',
        isTraversed && 'traversed',
    ].filter(Boolean).join(' ');

    return (
        <path
            className={classes}
            d={path}
            markerEnd="url(#arrowhead)"
            data-edge-id={id}
        />
    );
}

/**
 * EdgeMarkers - SVG defs for arrow markers.
 */
export function EdgeMarkers() {
    return (
        <defs>
            {/* Default arrowhead */}
            <marker
                id="arrowhead"
                markerWidth="10"
                markerHeight="7"
                refX="9"
                refY="3.5"
                orient="auto"
                markerUnits="strokeWidth"
            >
                <polygon
                    points="0 0, 10 3.5, 0 7"
                    fill="currentColor"
                />
            </marker>

            {/* Active arrowhead (glowing) */}
            <marker
                id="arrowhead-active"
                markerWidth="10"
                markerHeight="7"
                refX="9"
                refY="3.5"
                orient="auto"
                markerUnits="strokeWidth"
            >
                <polygon
                    points="0 0, 10 3.5, 0 7"
                    fill="var(--edge-active)"
                />
            </marker>

            {/* Glow filter for active elements */}
            <filter id="glow" x="-50%" y="-50%" width="200%" height="200%">
                <feGaussianBlur stdDeviation="3" result="coloredBlur" />
                <feMerge>
                    <feMergeNode in="coloredBlur" />
                    <feMergeNode in="SourceGraphic" />
                </feMerge>
            </filter>
        </defs>
    );
}
