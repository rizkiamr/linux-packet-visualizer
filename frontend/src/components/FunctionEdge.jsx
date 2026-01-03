/**
 * FunctionEdge - Renders a directed edge between two function nodes.
 */
export function FunctionEdge({
    id,
    fromX,
    fromY,
    toX,
    toY,
    orientation = 'vertical',
    isActive = false,
    isTraversed = false
}) {
    let path;

    if (orientation === 'horizontal') {
        const midX = (fromX + toX) / 2;
        path = `M ${fromX} ${fromY} C ${midX} ${fromY}, ${midX} ${toY}, ${toX} ${toY}`;
    } else {
        const midY = (fromY + toY) / 2;
        path = `M ${fromX} ${fromY} C ${fromX} ${midY}, ${toX} ${midY}, ${toX} ${toY}`;
    }

    const classes = [
        'function-edge',
        isActive && 'active',
        isTraversed && 'traversed',
    ].filter(Boolean).join(' ');

    return (
        <path
            className={classes}
            d={path}
            markerEnd={isActive ? "url(#arrowhead-active)" : "url(#arrowhead)"}
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
            {/* Default arrowhead (Open V-shape) */}
            <marker
                id="arrowhead"
                markerWidth="12"
                markerHeight="12"
                refX="10"
                refY="6"
                orient="auto"
                markerUnits="userSpaceOnUse"
            >
                <polyline
                    points="0 0, 10 6, 0 12"
                    fill="none"
                    stroke="var(--edge-default)"
                    strokeWidth="1.5"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                />
            </marker>

            {/* Active arrowhead (Open V-shape, glowing) */}
            <marker
                id="arrowhead-active"
                markerWidth="12"
                markerHeight="12"
                refX="10"
                refY="6"
                orient="auto"
                markerUnits="userSpaceOnUse"
            >
                <polyline
                    points="0 0, 10 6, 0 12"
                    fill="none"
                    stroke="var(--edge-active)"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
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
