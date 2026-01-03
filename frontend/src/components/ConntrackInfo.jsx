/**
 * ConntrackInfo - Displays current connection tracking state.
 */
export function ConntrackInfo({ conntrackState }) {
    if (!conntrackState) {
        return null;
    }

    // State colors
    const stateColors = {
        'NEW': '#ffe66d',
        'SYN_SENT': '#ff6b6b',
        'SYN_RECV': '#ff6b6b',
        'ESTABLISHED': '#00ff88',
        'FIN_WAIT': '#ffa500',
        'CLOSE_WAIT': '#ffa500',
        'LAST_ACK': '#ffa500',
        'TIME_WAIT': '#aaa',
        'CLOSED': '#666',
    };

    return (
        <div className="sidebar-section conntrack-info">
            <h2>Connection Tracking</h2>

            <div className="conntrack-state-container">
                <div
                    className="conntrack-state-badge"
                    style={{
                        backgroundColor: stateColors[conntrackState.state] || '#888',
                        color: conntrackState.state === 'ESTABLISHED' ? '#000' : '#fff'
                    }}
                >
                    {conntrackState.state}
                </div>
            </div>

            <p className="conntrack-description">
                {conntrackState.description}
            </p>

            <div className="conntrack-state-diagram">
                <div className="state-machine">
                    <div className={`state-node ${conntrackState.state === 'NEW' ? 'active' : ''}`}>NEW</div>
                    <div className="state-arrow">→</div>
                    <div className={`state-node ${conntrackState.state === 'SYN_SENT' ? 'active' : ''}`}>SYN_SENT</div>
                    <div className="state-arrow">→</div>
                    <div className={`state-node ${conntrackState.state === 'ESTABLISHED' ? 'active' : ''}`}>ESTABLISHED</div>
                </div>
            </div>
        </div>
    );
}
