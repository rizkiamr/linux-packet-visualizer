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

    // Mobile Breakpoint
    const isMobile = windowSize.width < 768;

    // SVG Layout Constants (Base values)
    const nodeHeight = isMobile ? 40 : 50;
    const nodeSpacingY = isMobile ? 60 : 80;
    const nodeGapX = isMobile ? 30 : 80; // Horizontal gap between nodes

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

    // Helper to calculate dynamic width
    const getDynamicWidth = (name) => {
        const charWidth = isMobile ? 7.5 : 9; // Approx px per char
        const padding = 40;
        const minW = isMobile ? 120 : 150;
        const textWidth = (name.length * charWidth) + padding;
        return Math.max(minW, textWidth);
    };

    // Calculate Layout
    const layout = useMemo(() => {
        if (!path?.functions) return null;

        const positions = {};

        const sidebarWidth = isMobile ? 0 : 320;
        const canvasWidth = windowSize.width - sidebarWidth - 40;
        const maxRowWidth = Math.max(isMobile ? 300 : 800, canvasWidth);

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
            if (funcs.length === 0) return;

            let row = 0;
            let currentX = 60;
            const layerStartY = currentY;

            funcs.forEach((fn) => {
                const thisWidth = getDynamicWidth(fn.name);

                // Check wrap
                if (currentX + thisWidth + 20 > maxRowWidth) {
                    row++;
                    currentX = 60;
                }

                // Layout: Grid with slight vertical separation for rows
                const y = layerStartY + 40 + (row * nodeSpacingY);

                positions[fn.id] = { x: currentX, y, width: thisWidth };
                currentX += thisWidth + nodeGapX;
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
            const sameRow = Math.abs(fromPos.y - toPos.y) < 20;

            // If same row and moving right -> Horizontal
            if (sameRow && toPos.x > fromPos.x) {
                return {
                    id: `${edge.from}-${edge.to}`,
                    fromX: fromPos.x + fromPos.width, // Use dynamic width
                    fromY: fromPos.y + nodeHeight / 2,
                    toX: toPos.x,
                    toY: toPos.y + nodeHeight / 2,
                    orientation: 'horizontal'
                };
            }

            // Vertical
            return {
                id: `${edge.from}-${edge.to}`,
                fromX: fromPos.x + fromPos.width / 2, // Center of dynamic width
                fromY: fromPos.y + nodeHeight,
                toX: toPos.x + toPos.width / 2, // Center of dynamic width
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
            x: pos.x + pos.width / 2,
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
                            width={pos.width} // Dynamic width
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
