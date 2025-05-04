<script lang="ts">
	import { mount } from "svelte";
	import { SearchAnime } from "$wails/go/main/App.js";
	import Anime from "$components/Anime.svelte";
	import { writable } from "svelte/store";

	let search: string;
	let enableSearch = writable(true);

	function doSearch(): void {
		enableSearch.set(false);
		SearchAnime(search).then((result) => {
			document.getElementById("result")!.innerHTML = ""
			result.forEach((element) =>
				mount(Anime, {
					props: { anime: element },
					target: document.getElementById("result")!,
				})
			);
			enableSearch.set(true);
		});
	}
</script>

<main>
	<form on:submit={(e) => {e.preventDefault(); doSearch()}}>
		<input
			autocomplete="off"
			bind:value={search}
			class="input"
			id="name"
			type="text"
		/>
		<button type="submit" disabled={!$enableSearch}>Cerca</button>
	</form>
	<div id="result"></div>
</main>

<style>
	#result {
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
	}
	main {
		display: flex;
		flex-direction: column;
		align-items: center;
	}
</style>
