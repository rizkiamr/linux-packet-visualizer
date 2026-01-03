/**
 * PathSelector - Dropdown to switch between egress and ingress paths.
 */
export function PathSelector({ paths, selectedPathId, onSelect }) {
    if (!paths || paths.length <= 1) {
        return null;
    }

    return (
        <div className="path-selector">
            <select
                value={selectedPathId || ''}
                onChange={(e) => onSelect(e.target.value)}
                className="path-select"
            >
                {paths.map((p) => (
                    <option key={p.id} value={p.id}>
                        {p.name}
                    </option>
                ))}
            </select>
            <div className="direction-badge" data-direction={paths.find(p => p.id === selectedPathId)?.direction}>
                {paths.find(p => p.id === selectedPathId)?.direction === 'egress' ? '↓ TX' : '↑ RX'}
            </div>
        </div>
    );
}
