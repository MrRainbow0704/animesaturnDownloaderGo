<script lang="ts">
import { onMount } from "svelte";
import {
	SearchAnime,
	GetDefaultAnime,
	GetPageNumber,
} from "$wails/go/main/App";
import { BrowserOpenURL } from "$wails/runtime";
import type { helper } from "$wails/go/models";
import Anime from "$components/Anime.svelte";
import SpinningWheeel from "$components/SpinningWheeel.svelte";
import ProgressBar from "$components/ProgressBar.svelte";
import Toast from "$components/Toast.svelte";

let search: string = $state("");
let mainText: string = $state("");
let load: boolean = $state(false);
let enableSearch = $state(true);
let anime: helper.Anime[] = $state([]);
let pages: number = $state(0);
let page: number = $state(1);
function defaultSearch(): void {
	pages = 0;
	mainText = "Anime in evidenza";
	load = true;
	enableSearch = false;
	GetDefaultAnime().then((result) => {
		anime = result;
		load = false;
		enableSearch = true;
	});
}
function doSearch({ noReset = false }: { noReset?: boolean } = {}): void {
	if (!search) {
		return defaultSearch();
	}
	mainText = `Risultati della ricerca "${search}"`;
	load = true;
	enableSearch = false;
	SearchAnime(search, page).then((result) => {
		anime = result;
		load = false;
		enableSearch = true;
	});
	GetPageNumber(search).then((result) => {
		pages = result;
		if (noReset === false) {
			console.log("Changing page")
			page = 1;
		}
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
				onclick={(e) => {
					e.preventDefault();
					BrowserOpenURL((e.target! as HTMLAnchorElement).href);
				}}>Marco Simone</a>
		</p>
	</header>
	<form
		onsubmit={(e) => {
			e.preventDefault();
			doSearch();
		}}>
		<input
			bind:value={search}
			required
			aria-required="true"
			autocomplete="off"
			spellcheck="false"
			class="input"
			id="name"
			type="text"
			placeholder="Cerca un'anime da scaricare" />
		<button type="submit" disabled={!enableSearch}>Cerca</button>
	</form>
	<div class="wrapper">
		<h2>{mainText}</h2>
		<div class="wrapper">
			{#if load}
				<SpinningWheeel />
			{:else if anime.length}
				<div id="result">
					{#each anime as a}
						<Anime anime={a} />
					{/each}
				</div>
				<nav>
					<ul>
						{#each { length: pages } as _, i}
							<li>
								<button
									type="button"
									disabled={i + 1 === page}
									aria-current={i + 1 === page
										? "page"
										: undefined}
									aria-disabled={i + 1 === page
										? "true"
										: "false"}
									onclick={(e) => {
										e.preventDefault();
										page = i + 1;
										doSearch({ noReset: false });
									}}>
									{i + 1}
								</button>
							</li>
						{/each}
					</ul>
				</nav>
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
.wrapper {
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
	width: 100%;
}
main {
	display: flex;
	flex-direction: column;
	align-items: center;
	width: 100%;
	height: 100%;
}
nav {
	width: 100%;
	display: flex;
	align-items: center;
	justify-content: center;
}
ul {
	display: flex;
	direction: row;
	flex-wrap: wrap;
	align-items: center;
	justify-content: center;
	list-style: none;
	padding: 0;
	max-width: 80%;
}
</style>
