/**
 * FunctionInfo - Sidebar panel showing details of the current function.
 */
export function FunctionInfo({ fn, kernelVersion }) {
    if (!fn) {
        return (
            <div className="sidebar-section function-info">
                <h2>Current Function</h2>
                <p style={{ color: 'var(--text-muted)', fontSize: '0.875rem' }}>
                    Press play or step to begin
                </p>
            </div>
        );
    }

    // Netfilter hook colors
    const nfHookColors = {
        'PREROUTING': '#ff6b6b',
        'INPUT': '#4ecdc4',
        'OUTPUT': '#ffe66d',
        'POSTROUTING': '#95e1d3',
        'FORWARD': '#a8e6cf',
    };

    // BPF hook colors
    const bpfHookColors = {
        'XDP': '#ff00ff',
        'TC_INGRESS': '#00bfff',
        'TC_EGRESS': '#1e90ff',
        'CGROUP_SKB': '#9370db',
        'SOCKET': '#da70d6',
    };

    // Generate kernel source URL using git.kernel.org
    const getKernelSourceUrl = () => {
        if (!fn.sourceFile) return null;
        const version = kernelVersion || '5.10.8';
        // Use exact version tag (e.g., "v5.10.8") to ensure line numbers match
        const tag = `v${version}`;
        const line = fn.lineNumber ? `#n${fn.lineNumber}` : '';
        return `https://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git/tree/${fn.sourceFile}?h=${tag}${line}`;
    };

    const sourceUrl = getKernelSourceUrl();

    return (
        <div className="sidebar-section function-info">
            <h2>Current Function</h2>
            <div className="fn-name">{fn.name}</div>
            {sourceUrl ? (
                <a
                    href={sourceUrl}
                    className="fn-file fn-file-link"
                    target="_blank"
                    rel="noopener noreferrer"
                    title="View source in Linux kernel"
                >
                    {fn.sourceFile}:{fn.lineNumber || '?'} ↗
                </a>
            ) : (
                <div className="fn-file">{fn.sourceFile}:{fn.lineNumber || '?'}</div>
            )}
            <div className="fn-description">{fn.description}</div>

            {fn.skbMutation && (
                <div className="mutation-badge">
                    <span>{fn.skbMutation.operation.toUpperCase()}</span>
                    <span>
                        {fn.skbMutation.headerType
                            ? `${fn.skbMutation.headerType} (${fn.skbMutation.size}B)`
                            : `${fn.skbMutation.size}B`
                        }
                    </span>
                </div>
            )}

            {fn.netfilterHook && (
                <div className="netfilter-info">
                    <h3>Netfilter Hook</h3>
                    <div
                        className="hook-badge"
                        style={{
                            backgroundColor: nfHookColors[fn.netfilterHook.hook] || '#888',
                            color: '#000'
                        }}
                    >
                        {fn.netfilterHook.hook}
                    </div>
                    <p className="hook-description">{fn.netfilterHook.description}</p>
                    <div className="hook-tables">
                        <span className="tables-label">Tables:</span>
                        {fn.netfilterHook.tables?.map((table, i) => (
                            <span key={i} className="table-tag">{table}</span>
                        ))}
                    </div>
                </div>
            )}

            {fn.bpfHook && (
                <div className="bpf-info">
                    <h3>eBPF Hook</h3>
                    <div
                        className="hook-badge bpf"
                        style={{
                            backgroundColor: bpfHookColors[fn.bpfHook.type] || '#888',
                            color: '#fff'
                        }}
                    >
                        {fn.bpfHook.type.replace('_', ' ')}
                    </div>
                    <p className="hook-description">{fn.bpfHook.description}</p>
                    <div className="hook-actions">
                        <span className="actions-label">Actions:</span>
                        {fn.bpfHook.actions?.map((action, i) => (
                            <span key={i} className="action-tag">{action}</span>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}
