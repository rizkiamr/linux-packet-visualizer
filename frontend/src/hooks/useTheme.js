import { useState, useEffect } from 'react';

/**
 * Hook to manage theme state (dark/light mode).
 * Persists preference in localStorage.
 */
export function useTheme() {
    const [theme, setTheme] = useState(() => {
        // Check localStorage first
        const saved = localStorage.getItem('theme');
        if (saved) return saved;

        // Check system preference
        if (window.matchMedia?.('(prefers-color-scheme: light)').matches) {
            return 'light';
        }

        return 'dark';
    });

    useEffect(() => {
        // Apply theme to document
        document.documentElement.setAttribute('data-theme', theme);

        // Persist to localStorage
        localStorage.setItem('theme', theme);
    }, [theme]);

    const toggle = () => {
        setTheme(t => t === 'dark' ? 'light' : 'dark');
    };

    const setDark = () => setTheme('dark');
    const setLight = () => setTheme('light');

    return {
        theme,
        isDark: theme === 'dark',
        isLight: theme === 'light',
        toggle,
        setDark,
        setLight,
    };
}
