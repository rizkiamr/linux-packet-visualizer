import { useState, useEffect, useCallback, useRef } from 'react';

/**
 * Hook for step-through simulation control.
 * @param {Array} simulation - Array of SimulateStep objects
 * @param {number} speed - Milliseconds between auto-advance steps
 * @returns Simulation state and controls
 */
export function useSimulation(simulation, speed = 1000) {
    const [stepIndex, setStepIndex] = useState(0);
    const [isPlaying, setIsPlaying] = useState(false);
    const [playSpeed, setPlaySpeed] = useState(speed);
    const intervalRef = useRef(null);

    const totalSteps = simulation?.length ?? 0;
    const currentStep = simulation?.[stepIndex] ?? null;
    const hasNext = stepIndex < totalSteps - 1;
    const hasPrev = stepIndex > 0;

    // Step forward
    const step = useCallback(() => {
        if (hasNext) {
            setStepIndex((i) => i + 1);
        } else {
            setIsPlaying(false);
        }
    }, [hasNext]);

    // Step backward
    const back = useCallback(() => {
        if (hasPrev) {
            setStepIndex((i) => i - 1);
        }
    }, [hasPrev]);

    // Reset to beginning
    const reset = useCallback(() => {
        setStepIndex(0);
        setIsPlaying(false);
    }, []);

    // Start auto-advance
    const play = useCallback(() => {
        if (hasNext) {
            setIsPlaying(true);
        }
    }, [hasNext]);

    // Stop auto-advance
    const pause = useCallback(() => {
        setIsPlaying(false);
    }, []);

    // Toggle play/pause
    const toggle = useCallback(() => {
        if (isPlaying) {
            pause();
        } else {
            play();
        }
    }, [isPlaying, play, pause]);

    // Jump to specific step
    const goTo = useCallback((index) => {
        if (index >= 0 && index < totalSteps) {
            setStepIndex(index);
        }
    }, [totalSteps]);

    // Set playback speed
    const setSpeed = useCallback((ms) => {
        setPlaySpeed(Math.max(100, Math.min(5000, ms)));
    }, []);

    // Auto-advance effect
    useEffect(() => {
        if (isPlaying && hasNext) {
            intervalRef.current = setInterval(() => {
                setStepIndex((i) => {
                    if (i >= totalSteps - 1) {
                        setIsPlaying(false);
                        return i;
                    }
                    return i + 1;
                });
            }, playSpeed);
        }

        return () => {
            if (intervalRef.current) {
                clearInterval(intervalRef.current);
            }
        };
    }, [isPlaying, hasNext, totalSteps, playSpeed]);

    // Compute traversed edges (edges taken before current step)
    const traversedEdges = simulation
        ?.slice(0, stepIndex)
        .filter((s) => s.edgeTaken)
        .map((s) => `${s.edgeTaken.from}-${s.edgeTaken.to}`) ?? [];

    return {
        currentStep,
        stepIndex,
        totalSteps,
        isPlaying,
        hasNext,
        hasPrev,
        traversedEdges,
        step,
        back,
        reset,
        play,
        pause,
        toggle,
        goTo,
        setSpeed,
        playSpeed,
    };
}
