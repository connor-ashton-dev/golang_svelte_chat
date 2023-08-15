import { writable } from 'svelte/store';

export const name = writable<string>('');
export const room = writable<string>('');
