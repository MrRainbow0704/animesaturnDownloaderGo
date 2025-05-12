<script lang="ts">
import { modalID } from "$lib/store";
import { getContext, type Snippet } from "svelte";

const modalId: string = getContext("modalId");
let { children }: { children?: Snippet } = $props();
</script>

<button aria-haspopup="dialog" onclick={() => ($modalID = modalId)}>
	{@render children?.()}
</button>

<style>
button {
	background: transparent;
	border: none;
}
button:hover {
	color: #fff;
	cursor: pointer;
	background-color: var(--secondary);
}
</style>
