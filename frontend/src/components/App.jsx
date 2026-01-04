import { usePacketPath } from '../hooks/usePacketPath';
import { useSimulation } from '../hooks/useSimulation';
import { useTheme } from '../hooks/useTheme';
import { SVGCanvas } from './SVGCanvas';
import { SKBuffDiagram } from './SKBuffDiagram';
import { FunctionInfo } from './FunctionInfo';
import { ConntrackInfo } from './ConntrackInfo';
import { Controls } from './Controls';
import { PathSelector } from './PathSelector';
import { ThemeToggle } from './ThemeToggle';

/**
 * App - Main application component.
 */
export function App() {
    const {
        paths,
        selectedPathId,
        setSelectedPath,
        path,
        simulation,
        metadata,
        kernelVersion,
        loading,
        error
    } = usePacketPath();

    const sim = useSimulation(simulation, 800);
    const { theme, toggle: toggleTheme } = useTheme();

    // Reset simulation when path changes
    const handlePathChange = (pathId) => {
        sim.reset();
        setSelectedPath(pathId);
    };

    // Jump to specific step when node is clicked
    const handleNodeClick = (functionId) => {
        const index = simulation?.findIndex(step => step.functionId === functionId || step.function.id === functionId);
        if (index !== -1) {
            sim.goTo(index);
        }
    };

    if (loading) {
        return (
            <div className="app">
                <div className="loading" style={{ gridColumn: '1 / -1', gridRow: '1 / -1' }}>
                    <div className="loading-spinner" />
                    <div className="loading-text">Loading Linux kernel packet path...</div>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="app">
                <div className="loading" style={{ gridColumn: '1 / -1', gridRow: '1 / -1' }}>
                    <div className="loading-text" style={{ color: 'var(--layer-user-accent)' }}>
                        Error: {error}
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="app">
            {/* Header */}
            <header className="app-header">
                <h1>The <span>Linux Packet Path</span></h1>
                {kernelVersion && (
                    <span className="kernel-version">Linux {kernelVersion}</span>
                )}
                <PathSelector
                    paths={paths}
                    selectedPathId={selectedPathId}
                    onSelect={handlePathChange}
                />
                <ThemeToggle theme={theme} onToggle={toggleTheme} />
            </header>

            {/* Main canvas */}
            <div className="canvas-container">
                <SVGCanvas
                    path={path}
                    currentFunctionId={sim.currentStep?.function?.id}
                    traversedEdges={sim.traversedEdges}
                    isPlaying={sim.isPlaying}
                    onNodeClick={handleNodeClick}
                />
            </div>

            {/* Sidebar */}
            <aside className="sidebar">
                <FunctionInfo fn={sim.currentStep?.function} kernelVersion={kernelVersion} />
                <ConntrackInfo conntrackState={sim.currentStep?.conntrackState} />
                <SKBuffDiagram
                    skbState={sim.currentStep?.skbuffState}
                    metadata={metadata}
                    direction={path?.direction}
                />
            </aside>

            {/* Controls */}
            <Controls
                stepIndex={sim.stepIndex}
                totalSteps={sim.totalSteps}
                isPlaying={sim.isPlaying}
                hasNext={sim.hasNext}
                hasPrev={sim.hasPrev}
                onStep={sim.step}
                onBack={sim.back}
                onToggle={sim.toggle}
                onReset={sim.reset}
            />
        </div>
    );
}
