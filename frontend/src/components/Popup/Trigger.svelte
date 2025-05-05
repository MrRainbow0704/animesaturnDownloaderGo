<script lang="ts">
	import { id } from "$lib/store";
	import { getContext } from "svelte";
	
	const modalId: string = getContext("modalId");
</script>

<button on:click={() => ($id = modalId)}>
	<slot />
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
