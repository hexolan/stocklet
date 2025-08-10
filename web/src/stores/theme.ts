import { writable } from 'svelte/store';

function getInitialDark(): boolean {
	if (typeof window === 'undefined') return false;
	try {
		const saved = localStorage.getItem('color-scheme');
		if (saved === 'dark') return true;
		if (saved === 'light') return false;
		return window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
	} catch {
		return false;
	}
}

export const darkMode = writable<boolean>(getInitialDark());

export function setDarkMode(value: boolean): void {
	darkMode.set(value);
}

export function toggleDarkMode(): void {
	darkMode.update((v) => !v);
}

if (typeof window !== 'undefined') {
	darkMode.subscribe((isDark) => {
		const root = document.documentElement;
		root.classList.toggle('dark', isDark);
		try {
			localStorage.setItem('color-scheme', isDark ? 'dark' : 'light');
		} catch {
			// ignore storage write errors
		}
	});
} 