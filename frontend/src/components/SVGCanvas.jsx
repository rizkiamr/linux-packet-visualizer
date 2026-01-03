import { useMemo, useState, useEffect } from 'react';
import { FunctionNode } from './FunctionNode';
import { FunctionEdge, EdgeMarkers } from './FunctionEdge';
import { PacketMarker } from './PacketMarker';

/**
 * SVGCanvas - Main visualization canvas with function nodes and edges.
 * Supports both egress (top-to-bottom) and ingress (bottom-to-top) layouts.
 */
export function SVGCanvas({
    path,
    currentFunctionId,
    traversedEdges = [],
    isPlaying = false,
    onNodeClick
}) {
    // Window size for dynamic layout
    const [windowSize, setWindowSize] = useState({ width: window.innerWidth, height: window.innerHeight });

    useEffect(() => {
        const handleResize = () => setWindowSize({ width: window.innerWidth, height: window.innerHeight });
        window.addEventListener('resize', handleResize);
        return () => window.removeEventListener('resize', handleResize);
    }, []);

    // SVG Layout Constants
    const nodeWidth = 150;
    const nodeHeight = 50;
    const nodeSpacingX = 250;
    const nodeSpacingY = 80; // Vertical spacing between rows in same layer

    // Determine direction
    const isIngress = path?.direction === 'ingress';

    // Layer Colors
    const layerColors = {
        'Device Driver': 'var(--layer-driver-bg)',
        'Data Link Layer': 'var(--layer-datalink-bg)',
        'Network Layer': 'var(--layer-network-bg)',
        'Transport Layer': 'var(--layer-transport-bg)',
        'Socket Layer': 'var(--layer-socket-bg)',
    };

    // Calculate Layout
    const layout = useMemo(() => {
        if (!path?.functions) return null;

        const positions = {};
        // Calculate available width for content (excluding sidebar approx 320px + margins)
        // Safer to assume some padding. 
        // We want to wrap BEFORE the horizontal scrollbar would appear.
        const canvasWidth = windowSize.width - 340;
        const maxRowWidth = Math.max(800, canvasWidth); // Ensure at least enough for a few nodes

        const layerOrder = isIngress
            ? ['Device Driver', 'Data Link Layer', 'Network Layer', 'Transport Layer', 'Socket Layer']
            : ['Transport Layer', 'Network Layer', 'Data Link Layer', 'Device Driver'];

        // Group functions
        const functionsByLayer = {};
        path.functions.forEach(fn => {
            if (!functionsByLayer[fn.layer]) functionsByLayer[fn.layer] = [];
            functionsByLayer[fn.layer].push(fn);
        });

        const layerBounds = [];
        let currentY = 50;

        layerOrder.forEach(layerName => {
            const funcs = functionsByLayer[layerName] || [];
            // Even if empty, we might want to show the layer? 
            // Current app only shows occupied layers usually, but let's stick to showing what we have.
            if (funcs.length === 0) return;

            let row = 0;
            let currentX = 60;
            const layerStartY = currentY;

            funcs.forEach((fn) => {
                // Check wrap
                if (currentX + nodeWidth + 20 > maxRowWidth) {
                    row++;
                    currentX = 60;
                }

                // Layout: Grid with slight vertical separation for rows
                const y = layerStartY + 40 + (row * nodeSpacingY);

                positions[fn.id] = { x: currentX, y };
                currentX += nodeSpacingX;
            });

            const layerHeight = (row + 1) * nodeSpacingY + 80;

            // Helper for class name
            let className = '';
            if (layerName.includes('Driver')) className = 'layer-driver';
            else if (layerName.includes('Data Link')) className = 'layer-datalink';
            else if (layerName.includes('Network')) className = 'layer-network';
            else if (layerName.includes('Transport')) className = 'layer-transport';
            else if (layerName.includes('Socket')) className = 'layer-socket';

            layerBounds.push({
                name: layerName,
                className,
                y: layerStartY,
                height: layerHeight,
                color: layerColors[layerName]
            });

            currentY += layerHeight + 20; // Gap between layers
        });

        return { positions, layerBounds, totalHeight: currentY + 100 };
    }, [path, isIngress, windowSize.width]);

    // Calculate edge positions
    const edgeData = useMemo(() => {
        if (!path?.edges || !layout) return [];

        return path.edges.map((edge) => {
            const fromPos = layout.positions[edge.from];
            const toPos = layout.positions[edge.to];

            if (!fromPos || !toPos) return null;

            // Determine if horizontal (same row)
            // Ideally checked by Y difference being small
            const sameRow = Math.abs(fromPos.y - toPos.y) < 20;

            // If same row and moving right -> Horizontal
            if (sameRow && toPos.x > fromPos.x) {
                return {
                    id: `${edge.from}-${edge.to}`,
                    fromX: fromPos.x + nodeWidth,
                    fromY: fromPos.y + nodeHeight / 2,
                    toX: toPos.x,
                    toY: toPos.y + nodeHeight / 2,
                    orientation: 'horizontal'
                };
            }

            // Otherwise (Cross-row or Cross-layer) -> Vertical
            // Right-to-Left or wrapping usually flows downwards or across.
            // Vertical connection points (Bottom -> Top) work best for wrapping too.
            return {
                id: `${edge.from}-${edge.to}`,
                fromX: fromPos.x + nodeWidth / 2,
                fromY: fromPos.y + nodeHeight,
                toX: toPos.x + nodeWidth / 2,
                toY: toPos.y,
                orientation: 'vertical'
            };
        }).filter(Boolean);
    }, [path?.edges, layout]);

    // Get current packet position
    const packetPos = useMemo(() => {
        if (!currentFunctionId || !layout?.positions[currentFunctionId]) {
            return null;
        }
        const pos = layout.positions[currentFunctionId];
        return {
            x: pos.x + nodeWidth / 2,
            y: pos.y + nodeHeight / 2,
        };
    }, [currentFunctionId, layout]);

    if (!path || !layout) {
        return (
            <div className="loading">
                <div className="loading-spinner" />
                <div className="loading-text">Loading kernel path data...</div>
            </div>
        );
    }

    return (
        <svg
            className={`svg-canvas ${isIngress ? 'ingress' : 'egress'}`}
            viewBox={`0 0 ${Math.max(windowSize.width - 320, 800)} ${layout.totalHeight}`}
            width="100%"
            height={layout.totalHeight}
            style={{ display: 'block', minWidth: '100%' }}
        >
            <EdgeMarkers />

            {/* Layer backgrounds */}
            <g className="layer-backgrounds">
                {layout.layerBounds.map((layer, i) => (
                    <rect
                        key={i}
                        className={`layer-bg ${layer.className}`}
                        x={0}
                        y={layer.y}
                        width="100%"
                        height={layer.height}
                        opacity={0.5}
                    />
                ))}
            </g>

            {/* Layer labels */}
            <g className="layer-labels">
                {layout.layerBounds.map((layer, i) => (
                    <text key={i} className="layer-label" x={20} y={layer.y + 20}>
                        {layer.name}
                    </text>
                ))}
            </g>

            {/* Edges */}
            <g className="edges">
                {edgeData.map((edge) => (
                    <FunctionEdge
                        key={edge.id}
                        {...edge}
                        isActive={
                            traversedEdges.includes(edge.id) ||
                            (packetPos && edge.toX === packetPos.x)
                        }
                        isTraversed={traversedEdges.includes(edge.id)}
                    />
                ))}
            </g>

            {/* Nodes */}
            <g className="nodes">
                {path.functions.map((fn) => {
                    const pos = layout.positions[fn.id];
                    if (!pos) return null;

                    return (
                        <FunctionNode
                            key={fn.id}
                            {...fn}
                            x={pos.x}
                            y={pos.y}
                            width={nodeWidth}
                            height={nodeHeight}
                            isActive={fn.id === currentFunctionId}
                            onClick={() => onNodeClick && onNodeClick(fn.id)}
                        />
                    );
                })}
            </g>

            {/* Packet */}
            {packetPos && (
                <PacketMarker
                    x={packetPos.x}
                    y={packetPos.y}
                    isMoving={isPlaying}
                />
            )}
        </svg>
    );
}
