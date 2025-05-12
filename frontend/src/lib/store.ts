import { writable, type Writable } from "svelte/store";

export const modalID: Writable<string> = writable("");
export const downloading: Writable<boolean> = writable(false);
