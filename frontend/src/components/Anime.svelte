<script lang="ts">
import { helper } from "$wails/go/models";
import AnimeMore from "$components/AnimeMore.svelte";

let { anime }: { anime: helper.Anime } = $props();
let dialog: HTMLDialogElement;
let body: HTMLBodyElement;
let isOpen: boolean = $state(false);
let height: number;
function close(): void {
	isOpen = false;
	dialog.close();
	window.scrollTo({
		top: height,
		behavior: "instant",
	});
	body.style.removeProperty("--height");
}
function open(): void {
	isOpen = true;
	height = window.scrollY;
	body.style.setProperty("--height", `${height}px`);
	dialog.showModal();
}
function handleKeydown(e: KeyboardEvent): void {
	if (e.key === "Escape") {
		e.preventDefault();
		if (dialog.open) {
			close();
		}
	}
}
</script>

<svelte:window onkeydown={handleKeydown} />
<svelte:body bind:this={body} />
<dialog aria-modal="true" id={anime.Url} bind:this={dialog}>
	<div class="content">
		<button class="close" type="button" onclick={close}>
			<span aria-hidden="true">&times;</span>
		</button>
		{#if isOpen}
			<AnimeMore {anime} />
		{/if}
	</div>
</dialog>
<button class="open" aria-haspopup="dialog" onclick={open}>
	<div class="anime" style="background-image: url({anime.Poster});">
		<h2>
			{anime.Title}
		</h2>
	</div>
</button>

<style>
.anime {
	background: transparent;
	background-position: center;
	background-repeat: no-repeat;
	background-size: cover;
	border: none;
	color: #fff;
	display: grid;
	grid-template-rows: auto;
	grid-template-columns: auto;
	text-align: center;
	justify-items: center;
	align-items: center;
	margin: 0.25rem;
	cursor: pointer;
	width: 16rem;
	height: auto;
	aspect-ratio: 2/3;
	margin-bottom: 0.5rem;
	position: relative;
}
h2 {
	position: absolute;
	bottom: 0;
	text-wrap: auto;
	padding: 0.5rem 0;
	margin: 0;
	background: rgba(0, 0, 0, 0.75);
	width: 100%;
}
:global(#app):has(dialog[open]) {
	position: fixed;
	bottom: var(--height);
}
dialog[open] {
	display: flex;
}
dialog {
	background: none;
	position: fixed;
	z-index: 10;
	left: 0;
	top: 0;
	width: 100%;
	height: 100%;
	justify-content: center;
	border: 0;
	margin: 0;
	color: #fff;
}
dialog::backdrop {
	background: rgba(0, 0, 0, 0.4);
}
.content {
	position: relative;
	background: transparent;
	background-color: var(--primary);
	padding: 2rem;
	margin: auto;
	height: 80%;
	width: 90%;
	overflow-y: auto;
}
.close {
	position: absolute;
	background: transparent;
	top: 0;
	right: 0;
	color: #fff;
	font-size: 28px;
	font-weight: bold;
	cursor: pointer;
	border: none;
	aspect-ratio: 1;
}
.open {
	background: transparent;
	border: none;
}
.open:hover {
	color: #fff;
	cursor: pointer;
	background-color: var(--secondary);
}
</style>
