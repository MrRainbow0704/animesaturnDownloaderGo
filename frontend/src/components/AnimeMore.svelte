<script lang="ts">
	import type { helper } from "$wails/go/models";
	import { DownloadAnime } from "$wails/go/main/App";
	import { writable } from "svelte/store";

	export let anime: helper.Anime;

	let primo: number = 0;
	let ultimo: number = anime.Info.EpisodeCount;
	let filename: string = anime.Title.replace(/[<>:"/\\|?*]+/g,'');;
	let workers: number = 3;
	let downloadStatus = writable("");
	function download(): void {
		DownloadAnime(anime.Url, primo, ultimo, filename, workers).then((ok) =>
			ok ? downloadStatus.set("Fatto") : downloadStatus.set("Fallito")
		);
	}
</script>

<div class="more">
	<!-- <header>
		<img src={anime.Banner} alt="{anime.Title} banner" />
	</header> -->
	<aside>
		<img src={anime.Poster} alt="{anime.Title} poster" />
	</aside>
	<article>
		<h1>{anime.Title}</h1>
		<hr />
		<ul class="stats">
			<li><b>Studio:</b> {anime.Info.Studio}</li>
			<hr />
			<li><b>Status:</b> {anime.Info.Status}</li>
			<hr />
			<li><b>Episodi:</b> {anime.Info.EpisodeCount}</li>
			<hr />
			<li>
				<b>Tags:</b>
				<ul class="tags">
					{#each anime.Info.Tags as tag}
						<li>{tag}</li>
					{/each}
				</ul>
			</li>
		</ul>
		<hr />
		<h2>Trama</h2>
		<p>{anime.Info.Plot}</p>
		<hr />
		<h2>Download</h2>
		<form onsubmit={download}>
			<label>
				Nome dei File scaricati:
				<input type="text" name="filename" bind:value={filename} />
				[{primo}-{ultimo}].mp4
			</label>
			<div class="input-container">
				Episodi da scaricare.
				<label>
					Da:
					<input
						type="number"
						name="primo"
						min="0"
						max={anime.Info.EpisodeCount}
						bind:value={primo}
					/>
				</label>
				<label>
					a:
					<input
						type="number"
						name="ultimo"
						min="0"
						max={anime.Info.EpisodeCount}
						bind:value={ultimo}
					/>
				</label>
			</div>
			<label>
				Quanti worker da utilizzare
				<input
					type="number"
					name="workers"
					min="0"
					max="16"
					bind:value={workers}
				/>
			</label>
			<button type="submit">Download</button>
			<div>{$downloadStatus}</div>
		</form>
	</article>
</div>

<style>
	aside {
		padding-right: 1rem;
	}
	article {
		padding-left: 1rem;
	}
	aside > img {
		max-width: 100%;
	}
	.more {
		display: grid;
		grid-template-columns: 30% 70%;
	}
	ul {
		list-style: none;
		padding-inline-start: 0;
	}
	ul.tags {
		display: inline-block;
	}
	ul.tags > li {
		display: inline-block;
		background: #fff;
		color: var(--primary);
		padding: 0.25rem;
		margin: 0 0.125rem;
		border-radius: 4px;
	}
	ul > hr {
		border-color: var(--tertiary);
	}
	h2 {
		margin: 1rem 0 0.25rem 0;
	}
	p {
		margin: 0;
	}
	form > label {
		display: grid;
		grid-template-columns: 15rem auto;
		padding: 0.25rem 0;
	}
	form > label:has(input[type="number"]) {
		grid-template-columns: 15rem 5rem;
	}
	.input-container {
		display: grid;
		grid-template-columns: 15rem 10rem 10rem;
	}
	form > label:has(input[name="filename"]) {
		grid-template-columns: 15rem auto 6rem;
	}
</style>
