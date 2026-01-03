/**
 * PacketMarker - Animated indicator showing current packet position.
 */
export function PacketMarker({ x, y, isMoving = false }) {
    const classes = [
        'packet-marker',
        isMoving && 'moving',
    ].filter(Boolean).join(' ');

    return (
        <g className={classes} transform={`translate(${x}, ${y})`}>
            {/* Outer glow ring */}
            <circle r={16} fill="none" stroke="var(--packet-glow)" strokeWidth={2} opacity={0.3} />
            <circle r={12} fill="none" stroke="var(--packet-glow)" strokeWidth={2} opacity={0.5} />

            {/* Core */}
            <circle r={8} fill="var(--packet-glow)" />

            {/* Inner highlight */}
            <circle r={4} fill="white" opacity={0.6} />
        </g>
    );
}
