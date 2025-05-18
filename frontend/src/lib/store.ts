import { writable, type Writable } from "svelte/store";

export const downloading: Writable<boolean> = writable(false);
