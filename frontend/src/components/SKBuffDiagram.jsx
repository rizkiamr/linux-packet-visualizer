/**
 * SKBuffDiagram - Visual representation of sk_buff memory layout.
 */
export function SKBuffDiagram({ skbState, metadata }) {
    if (!skbState) {
        return (
            <div className="sidebar-section">
                <h2>sk_buff</h2>
                <div className="skbuff-diagram">
                    <p style={{ color: 'var(--text-muted)', fontSize: '0.875rem' }}>
                        No packet data
                    </p>
                </div>
            </div>
        );
    }

    const { head, data, tail, end, layers = [] } = skbState;
    const totalSize = end - head;
    const headroom = data - head;
    const dataLen = tail - data;
    const tailroom = end - tail;

    // Calculate percentages for visual bar
    const headroomPct = (headroom / totalSize) * 100;
    const dataPct = (dataLen / totalSize) * 100;
    const tailroomPct = (tailroom / totalSize) * 100;

    // Calculate header segments within data
    const headerSegments = layers.map((layer) => ({
        protocol: layer.protocol,
        size: layer.size,
        pct: (layer.size / totalSize) * 100,
    }));

    // Payload is data minus all headers
    const headerTotal = layers.reduce((sum, l) => sum + l.size, 0);
    const payloadSize = dataLen - headerTotal;
    const payloadPct = (payloadSize / totalSize) * 100;

    return (
        <div className="sidebar-section">
            <h2>sk_buff Memory Layout</h2>
            <div className="skbuff-diagram">
                {/* Visual bar */}
                <div className="skbuff-bar">
                    {/* Headroom */}
                    <div
                        className="skbuff-segment headroom"
                        style={{ flexBasis: `${headroomPct}%` }}
                        title={`Headroom: ${headroom} bytes`}
                    >
                        {headroomPct > 8 && 'headroom'}
                    </div>

                    {/* Protocol headers */}
                    {headerSegments.map((seg, i) => (
                        <div
                            key={i}
                            className={`skbuff-segment ${seg.protocol}`}
                            style={{ flexBasis: `${seg.pct}%` }}
                            title={`${seg.protocol.toUpperCase()}: ${seg.size} bytes`}
                        >
                            {seg.pct > 3 && seg.protocol.toUpperCase()}
                        </div>
                    ))}

                    {/* Payload */}
                    <div
                        className="skbuff-segment payload"
                        style={{ flexBasis: `${payloadPct}%` }}
                        title={`Payload: ${payloadSize} bytes`}
                    >
                        {payloadPct > 10 && 'PAYLOAD'}
                    </div>

                    {/* Tailroom */}
                    <div
                        className="skbuff-segment tailroom"
                        style={{ flexBasis: `${tailroomPct}%` }}
                        title={`Tailroom: ${tailroom} bytes`}
                    >
                        {tailroomPct > 8 && 'tailroom'}
                    </div>
                </div>

                {/* Pointer labels */}
                <div className="skbuff-pointers">
                    <span>HEAD ({head})</span>
                    <span style={{ marginLeft: `${headroomPct - 5}%` }}>DATA ({data})</span>
                    <span>TAIL ({tail})</span>
                    <span>END ({end})</span>
                </div>

                {/* Stats grid */}
                <div className="skbuff-stats">
                    <div className="skbuff-stat">
                        <div className="label">Packet Length</div>
                        <div className="value">{dataLen} bytes</div>
                    </div>
                    <div className="skbuff-stat">
                        <div className="label">Headroom</div>
                        <div className="value">{headroom} bytes</div>
                    </div>
                    <div className="skbuff-stat">
                        <div className="label">Headers</div>
                        <div className="value">{headerTotal} bytes</div>
                    </div>
                    <div className="skbuff-stat">
                        <div className="label">Payload</div>
                        <div className="value">{payloadSize} bytes</div>
                    </div>
                </div>
            </div>
        </div>
    );
}
