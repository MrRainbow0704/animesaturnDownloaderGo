<script lang="ts">
	import { downloading } from "$lib/store";
	import { writable, type Writable } from "svelte/store";

	const percentage: Writable<string> = writable("");
	let progressBar = setInterval(() => {
		// @ts-expect-error
		let progress: number = window.progressBarProgress;
		// @ts-expect-error
		let total: number = window.progressBarTotal;
		percentage.set(((progress / total) * 100).toFixed(2));

		if (!downloading) {
			clearInterval(progressBar);
		}
	}, 1000);
</script>

<footer>
	{#if $downloading}
		<div style="--text: '{$percentage}%';">
			<span style="width: {$percentage}%;"></span>
		</div>
	{/if}
</footer>

<style>
	div {
    position: relative;
		height: 100%;
		width: 100%;
		border-radius: 4px 4px 0 0;
		background: var(--secondary);
	}

	span {
    width: 0;
		display: block;
		background: var(--accent);
		height: 100%;
    border-radius: 0 1rem 1rem 0;
	}

	span::after {
		content: var(--text);
		text-wrap: nowrap;
		float: right;
		position: relative;
		right: 0.5rem;
		top: calc(0.5rem - 2px);
    z-index: 5;
	}

	footer {
		display: flex;
		position: fixed;
		z-index: 11;
		bottom: 0;
		left: 0;
		justify-content: center;
		height: 2rem;
		width: 100%;
	}
</style>
