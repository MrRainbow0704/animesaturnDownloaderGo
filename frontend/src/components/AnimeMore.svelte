<script lang="ts">
import { helper } from "$wails/go/models";
import { DownloadAnime, GetAnimeInfo } from "$wails/go/main/App";
import { downloading } from "$lib/store";
import { notifications } from "$lib/notifications";
import SpinningWheeel from "$src/components/SpinningWheeel.svelte";

let { anime }: { anime: helper.Anime } = $props();

let info: helper.AnimeInfo = $state(
	new helper.AnimeInfo({
		EpisodeCount: 0,
		Is18plus: false,
		Tags: [],
		Studio: "",
		Status: "",
		Plot: "",
		EpisodesList: [],
	})
);
let loaded: boolean = $state(false);
GetAnimeInfo(anime.Url).then((res) => {
	loaded = true;
	info = res;
});

let primo: string = $derived(info.EpisodesList[0]);
let ultimo: string = $derived(info.EpisodesList[info.EpisodesList.length - 1]);
let filename: string = $derived(anime.Title.replace(/[<>:"/\\|?*]+/g, ""));
let workers: number = $derived(info.EpisodeCount < 3 ? info.EpisodeCount : 3);
let downloadStatus: string = $state("");

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

{#if loaded}
	<div class="more">
		<aside aria-hidden="true">
			<img src={anime.Poster} alt="{anime.Title} poster" />
		</aside>
		<article>
			<h1>{anime.Title}</h1>
			<hr />
			{#if info.Is18plus}
				<div class="hentai">
					Stai per vedere una serie <b>R-18</b>.
					<br />
					Questa serie è adatta solo ad un pubblico
					<b>maggiorenne</b>.
				</div>
			{/if}
			<ul class="stats">
				<li><b>Studio:</b> {info.Studio}</li>
				<hr />
				<li><b>Status:</b> {info.Status}</li>
				<hr />
				<li>
					<b>Episodi:</b>
					{!info.EpisodeCount ? "??" : info.EpisodeCount}
				</li>
				<hr />
				<li>
					<b>Tags:</b>
					<ul class="tags">
						{#each info.Tags as tag}
							<li>{tag}</li>
						{/each}
					</ul>
				</li>
			</ul>
			<hr />
			<h2>Trama</h2>
			<p>{info.Plot}</p>
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
							{#each info.EpisodesList as e}
								<option
									selected={e === info.EpisodesList[0]}
									value={e}>{e}</option>
							{/each}
						</select>
						<label>
							a:
							<select name="ultimo" bind:value={ultimo}>
								{#each info.EpisodesList as e}
									<option
										selected={e ===
											info.EpisodesList[
												info.EpisodesList.length - 1
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
						max={info.EpisodeCount}
						bind:value={workers} />
				</label>
				<button disabled={$downloading} type="submit">Download</button>
				<div>{downloadStatus}</div>
			</form>
		</article>
	</div>
{:else}
	<SpinningWheeel />
{/if}

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
