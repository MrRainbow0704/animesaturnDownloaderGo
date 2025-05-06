<script lang="ts">
	import { flip } from "svelte/animate";
	import { fly } from "svelte/transition";
	import { notifications } from "$src/lib/notifications";

	let themes: {
		[key: string]: string | undefined;
		success: string;
		error: string;
		info: string;
		default: string;
	} = {
		success: "#3f6212",
		error: "#991b1b",
		info: "#06b6d4",
		default: "#434546",
	};
</script>

<div class="notifications">
	{#each $notifications as notification (notification.id)}
		<div
			animate:flip
			class="toast"
			style="background: {themes[notification.type]};"
			transition:fly={{ y: 30 }}
		>
			<div class="content">{notification.message}</div>
		</div>
	{/each}
</div>

<style>
	.notifications {
		position: fixed;
		top: 10px;
		left: 0;
		right: 0;
		margin: 0 auto;
		padding: 0;
		z-index: 12;
		display: flex;
		flex-direction: column;
		justify-content: flex-start;
		align-items: flex-end;
		pointer-events: none;
	}

	.toast {
		flex: 0 0 auto;
		margin-bottom: 0.5rem;
		margin-right: 1rem;
		border-radius: 4px;
		min-width: 14rem;
	}

	.content {
		padding: 10px;
		display: block;
		color: white;
		font-weight: 500;
	}
</style>
