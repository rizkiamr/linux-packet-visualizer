/**
 * Controls - Playback controls for the simulation.
 */
export function Controls({
    stepIndex,
    totalSteps,
    isPlaying,
    hasNext,
    hasPrev,
    onStep,
    onBack,
    onToggle,
    onReset
}) {
    return (
        <div className="controls-container">
            {/* Reset button */}
            <button
                className="control-btn"
                onClick={onReset}
                title="Reset to beginning"
                disabled={stepIndex === 0}
            >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={2}>
                    <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
                    <path d="M3 3v5h5" />
                </svg>
            </button>

            {/* Back button */}
            <button
                className="control-btn"
                onClick={onBack}
                title="Previous step"
                disabled={!hasPrev}
            >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={2}>
                    <polygon points="19 20 9 12 19 4 19 20" fill="currentColor" />
                    <line x1="5" y1="4" x2="5" y2="20" />
                </svg>
            </button>

            {/* Play/Pause button */}
            <button
                className="control-btn primary"
                onClick={onToggle}
                title={isPlaying ? 'Pause' : 'Play'}
                disabled={!hasNext && !isPlaying}
            >
                {isPlaying ? (
                    <svg viewBox="0 0 24 24" fill="currentColor">
                        <rect x="6" y="4" width="4" height="16" />
                        <rect x="14" y="4" width="4" height="16" />
                    </svg>
                ) : (
                    <svg viewBox="0 0 24 24" fill="currentColor">
                        <polygon points="5 3 19 12 5 21 5 3" />
                    </svg>
                )}
            </button>

            {/* Forward button */}
            <button
                className="control-btn"
                onClick={onStep}
                title="Next step"
                disabled={!hasNext}
            >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={2}>
                    <polygon points="5 4 15 12 5 20 5 4" fill="currentColor" />
                    <line x1="19" y1="4" x2="19" y2="20" />
                </svg>
            </button>

            {/* Step indicator */}
            <div className="step-indicator">
                Step <span className="current">{stepIndex + 1}</span> / {totalSteps}
            </div>
        </div>
    );
}
