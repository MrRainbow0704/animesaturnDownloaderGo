<script lang="ts">
import type { helper } from "$wails/go/models";
import { DownloadAnime } from "$wails/go/main/App";
import { downloading } from "$lib/store";
import { notifications } from "$lib/notifications";

let { anime }: { anime: helper.Anime } = $props();

let primo: string = $derived(anime.Info.EpisodesList[0]);
let ultimo: string = $derived(
	anime.Info.EpisodesList[anime.Info.EpisodesList.length - 1]
);
let filename: string = $derived(anime.Title.replace(/[<>:"/\\|?*]+/g, ""));
let workers: number = $derived(
	anime.Info.EpisodeCount < 3 ? anime.Info.EpisodeCount : 3
);
let downloadStatus = $state("");

function download(): void {
	if (primo > ultimo) {
		notifications.error(
			"Il primo episodio da scaricare non può essere prima dell'ultimo!",
			3000
		);
		return;
	}
	downloading.set(true);
	DownloadAnime(
		anime.Url,
		parseInt(primo),
		parseInt(ultimo),
		filename,
		workers
	).then((ok) => {
		downloading.set(false);
		ok
			? notifications.success("Finito di scaricare i file!", 3000)
			: notifications.error("Download fallito! :(", 3000);
	});
}
</script>

<div class="more">
	<aside aria-hidden="true">
		<img src={anime.Poster} alt="{anime.Title} poster" />
	</aside>
	<article>
		<h1>{anime.Title}</h1>
		<hr />
		{#if anime.Info.Is18plus}
			<div class="hentai">
				Stai per vedere una serie <b>R-18</b>.
				<br />
				Questa serie è adatta solo ad un pubblico <b>maggiorenne</b>.
			</div>
		{/if}
		<ul class="stats">
			<li><b>Studio:</b> {anime.Info.Studio}</li>
			<hr />
			<li><b>Status:</b> {anime.Info.Status}</li>
			<hr />
			<li>
				<b>Episodi:</b>
				{!anime.Info.EpisodeCount ? "??" : anime.Info.EpisodeCount}
			</li>
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
		<form
			onsubmit={(e) => {
				e.preventDefault();
				download();
			}}>
			<label style="--text: '[{primo}-{ultimo}].mp4';">
				Nome dei File scaricati:
				<input type="text" name="filename" bind:value={filename} />
			</label>
			<div class="input-container">
				Episodi da scaricare.
				<label>
					Da:
					<select name="primo" bind:value={primo}>
						{#each anime.Info.EpisodesList as e}
							<option
								selected={e === anime.Info.EpisodesList[0]}
								value={e}>{e}</option>
						{/each}
					</select>
					<label>
						a:
						<select name="ultimo" bind:value={ultimo}>
							{#each anime.Info.EpisodesList as e}
								<option
									selected={e ===
										anime.Info.EpisodesList[
											anime.Info.EpisodesList.length - 1
										]}
									value={e}>{e}</option>
							{/each}
						</select>
					</label>
				</label>
			</div>
			<label>
				Quanti worker da utilizzare
				<input
					type="number"
					name="workers"
					min="1"
					max={anime.Info.EpisodeCount}
					bind:value={workers} />
			</label>
			<button disabled={$downloading} type="submit">Download</button>
			<div>{downloadStatus}</div>
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
	width: 100%;
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
form > label:has(input[type="text"])::after {
	content: var(--text);
	text-wrap: nowrap;
	display: inline-block;
	position: relative;
	right: 5rem;
	top: 2px;
}
.input-container {
	display: grid;
	grid-template-columns: 15rem 10rem 10rem;
}
form > label:has(input[name="filename"]) {
	grid-template-columns: 15rem auto 0;
}
.hentai {
	position: relative;
	background: #ffc3c3;
	color: #bc2e2e;
	text-align: center;
	border-radius: 4px;
	width: 90%;
	margin: auto;
	padding: 0.5rem;
	font-size: 1.25rem;
}
.hentai::before,
.hentai::after {
	content: "⚠";
	font-size: 2.5rem;
	position: absolute;
	top: 0.25rem;
}
.hentai::before {
	left: 1rem;
}
.hentai::after {
	right: 1rem;
}
</style>
