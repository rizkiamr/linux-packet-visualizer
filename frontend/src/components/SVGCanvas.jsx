import { useMemo } from 'react';
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
    isPlaying = false
}) {
    // SVG dimensions
    const width = 1400;
    const height = 900;

    // Node dimensions
    const nodeWidth = 150;
    const nodeHeight = 50;
    const nodeSpacingX = 180;
    const nodeSpacingY = 70;

    // Determine direction
    const isIngress = path?.direction === 'ingress';

    // Layer configuration - different order for egress vs ingress
    const layerConfig = useMemo(() => {
        if (isIngress) {
            // Ingress: Driver at top, Socket at bottom (packet rises)
            return {
                'Device Driver': { baseY: 80, color: 'var(--layer-driver-bg)' },
                'Data Link Layer': { baseY: 280, color: 'var(--layer-datalink-bg)' },
                'Network Layer': { baseY: 480, color: 'var(--layer-network-bg)' },
                'Transport Layer': { baseY: 680, color: 'var(--layer-transport-bg)' },
                'Socket Layer': { baseY: 800, color: 'var(--layer-socket-bg)' },
            };
        } else {
            // Egress: Transport at top, Driver at bottom (packet descends)
            return {
                'Transport Layer': { baseY: 80, color: 'var(--layer-transport-bg)' },
                'Network Layer': { baseY: 320, color: 'var(--layer-network-bg)' },
                'Data Link Layer': { baseY: 560, color: 'var(--layer-datalink-bg)' },
                'Device Driver': { baseY: 760, color: 'var(--layer-driver-bg)' },
            };
        }
    }, [isIngress]);

    // Layer backgrounds configuration
    const layerBgs = useMemo(() => {
        if (isIngress) {
            return [
                { className: 'layer-driver', y: 50, height: 200 },
                { className: 'layer-datalink', y: 250, height: 200 },
                { className: 'layer-network', y: 450, height: 200 },
                { className: 'layer-transport', y: 650, height: 150 },
                { className: 'layer-socket', y: 800, height: 100 },
            ];
        } else {
            return [
                { className: 'layer-transport', y: 50, height: 240 },
                { className: 'layer-network', y: 290, height: 240 },
                { className: 'layer-datalink', y: 530, height: 200 },
                { className: 'layer-driver', y: 730, height: 170 },
            ];
        }
    }, [isIngress]);

    // Layer labels configuration
    const layerLabels = useMemo(() => {
        if (isIngress) {
            return [
                { text: 'Device Driver', y: 70 },
                { text: 'Data Link Layer (L2)', y: 270 },
                { text: 'Network Layer (L3)', y: 470 },
                { text: 'Transport Layer (L4)', y: 670 },
                { text: 'Socket Layer', y: 820 },
            ];
        } else {
            return [
                { text: 'Transport Layer (L4)', y: 70 },
                { text: 'Network Layer (L3)', y: 310 },
                { text: 'Data Link Layer (L2)', y: 550 },
                { text: 'Device Driver', y: 750 },
            ];
        }
    }, [isIngress]);

    // Calculate node positions
    const nodePositions = useMemo(() => {
        if (!path?.functions) return {};

        const positions = {};
        const layerCounts = {};

        // Group functions by layer
        path.functions.forEach((fn) => {
            const layer = fn.layer;
            if (!layerCounts[layer]) {
                layerCounts[layer] = 0;
            }

            const config = layerConfig[layer] || { baseY: 400 };
            const count = layerCounts[layer];

            // Calculate position
            const x = 60 + (count * nodeSpacingX);
            const y = config.baseY + (count % 2 === 0 ? 0 : nodeSpacingY);

            positions[fn.id] = { x, y };
            layerCounts[layer]++;
        });

        return positions;
    }, [path?.functions, layerConfig]);

    // Calculate edge positions
    const edgeData = useMemo(() => {
        if (!path?.edges) return [];

        return path.edges.map((edge) => {
            const fromPos = nodePositions[edge.from];
            const toPos = nodePositions[edge.to];

            if (!fromPos || !toPos) return null;

            return {
                id: `${edge.from}-${edge.to}`,
                fromX: fromPos.x + nodeWidth / 2,
                fromY: fromPos.y + nodeHeight,
                toX: toPos.x + nodeWidth / 2,
                toY: toPos.y,
            };
        }).filter(Boolean);
    }, [path?.edges, nodePositions]);

    // Get current packet position
    const packetPos = useMemo(() => {
        if (!currentFunctionId || !nodePositions[currentFunctionId]) {
            return null;
        }
        const pos = nodePositions[currentFunctionId];
        return {
            x: pos.x + nodeWidth / 2,
            y: pos.y + nodeHeight / 2,
        };
    }, [currentFunctionId, nodePositions]);

    if (!path) {
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
            viewBox={`0 0 ${width} ${height}`}
            preserveAspectRatio="xMidYMid meet"
        >
            <EdgeMarkers />

            {/* Layer backgrounds */}
            <g className="layer-backgrounds">
                {layerBgs.map((layer, i) => (
                    <rect
                        key={i}
                        className={`layer-bg ${layer.className}`}
                        x={0}
                        y={layer.y}
                        width={width}
                        height={layer.height}
                        opacity={0.5}
                    />
                ))}
            </g>

            {/* Layer labels */}
            <g className="layer-labels">
                {layerLabels.map((label, i) => (
                    <text key={i} className="layer-label" x={20} y={label.y}>
                        {label.text}
                    </text>
                ))}
            </g>

            {/* Edges (drawn first, behind nodes) */}
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

            {/* Function nodes */}
            <g className="nodes">
                {path.functions.map((fn) => {
                    const pos = nodePositions[fn.id];
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
                        />
                    );
                })}
            </g>

            {/* Packet marker */}
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
