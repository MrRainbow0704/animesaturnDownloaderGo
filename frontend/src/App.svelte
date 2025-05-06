<script lang="ts">
	import { onMount } from "svelte";
	import { writable } from "svelte/store";
	import { SearchAnime, GetDefaultAnime } from "$wails/go/main/App.js";
	import { BrowserOpenURL } from "$wails/runtime";
	import type { helper } from "$wails/go/models";
	import Anime from "$components/Anime.svelte";
	import SpinningWheeel from "$components/SpinningWheeel.svelte";
	import ProgressBar from "$components/ProgressBar.svelte";
	import Toast from "$components/Toast.svelte";

	let search: string;
	let mainText: string;
	let load: boolean = false;
	let enableSearch = writable(true);
	let anime: Array<helper.Anime> = [];
	function defaultSearch(): void {
		mainText = "Anime in evidenza";
		load = true;
		enableSearch.set(false);
		GetDefaultAnime().then((result) => {
			anime = result;
			load = false;
			enableSearch.set(true);
		});
	}
	function doSearch(): void {
		if (!search) {
			return defaultSearch();
		}
		mainText = `Risultati della ricerca "${search}"`;
		load = true;
		enableSearch.set(false);
		SearchAnime(search).then((result) => {
			anime = result;
			load = false;
			enableSearch.set(true);
		});
	}

	onMount(defaultSearch);
</script>

<main>
	<header>
		<h1>
			Anime<span style="color: var(--accent);">Saturn</span> Downlader
		</h1>
		<p>
			Creato da <a
				href="https://github.com/MrRainbow0704"
				on:click={(e) => {
					e.preventDefault();
					BrowserOpenURL((e.target! as HTMLAnchorElement).href);
				}}>Marco Simone</a
			>
		</p>
	</header>
	<form
		on:submit={(e) => {
			e.preventDefault();
			doSearch();
		}}
	>
		<input
			autocomplete="off"
			bind:value={search}
			class="input"
			id="name"
			type="text"
			placeholder="Cerca un'anime da scaricare"
		/>
		<button type="submit" disabled={!$enableSearch}>Cerca</button>
	</form>
	<div id="wrapper">
		<h2>{mainText}</h2>
		<div id="result">
			{#if load}
				<SpinningWheeel />
			{:else if anime.length}
				{#each anime as a}
					<Anime anime={a} />
				{/each}
			{:else}
				<p>Nessun risultato trovato</p>
			{/if}
		</div>
	</div>
	<ProgressBar />
	<Toast />
</main>

<style>
	h1,
	h2,
	p {
		margin: 0;
		padding: 0;
		text-align: center;
	}
	h2 {
		margin-bottom: 0.75rem;
	}
	a {
		color: var(--accent);
	}
	form {
		display: flex;
		flex-wrap: nowrap;
		flex-direction: row;
		justify-content: center;
		margin: 0.5rem;
		width: 100%;
	}
	input {
		font-size: 1rem;
		padding: 0.5rem;
		width: 60%;
	}
	#wrapper {
		display: flex;
		flex-wrap: nowrap;
		flex-direction: column;
		align-items: center;
		width: 100%;
		height: 100%;
	}
	#result {
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		height: 100%;
		width: 100%;
	}
	main {
		display: flex;
		flex-direction: column;
		align-items: center;
		width: 100%;
		height: 100%;
	}
</style>
