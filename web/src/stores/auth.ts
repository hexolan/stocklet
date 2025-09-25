import { writable } from 'svelte/store';
import type { components } from "$lib/api/schema";

type AuthStore = {
  profile: components["schemas"]["v1User"] | null;
  tokens: components["schemas"]["v1AuthToken"] | null;
  isLoading: boolean;
};

const INITIAL_AUTH_STATE: AuthStore = {
  profile: null,
  tokens: null,
  isLoading: false
};

export const authState = writable<AuthStore>(INITIAL_AUTH_STATE);

export function setTokens(value: components["schemas"]["v1AuthToken"]): void {
  authState.update((currentState) => {
    return { ...currentState, tokens: value }
  });
}

export function setProfile(value: components["schemas"]["v1User"]): void {
  authState.update((currentState) => {
    return { ...currentState, profile: value }
  });
}