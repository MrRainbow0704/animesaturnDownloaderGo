<script lang="ts">
	// import { mount } from "svelte";
	import { helper } from "$wails/go/models";
	import AnimeMore from "$components/AnimeMore.svelte";
	import Modal from "$components/Popup/Modal.svelte";
	import Content from "$components/Popup/Content.svelte";
	import Trigger from "$components/Popup/Trigger.svelte";

	export let anime: helper.Anime;
</script>

<Modal modalId={anime.Url}>
	<Content>
		<AnimeMore {anime} />
	</Content>
	<Trigger>
		<div>
			<span><img src={anime.Poster} alt="{anime.Title} psoter" /></span>
			<h2>{anime.Title}</h2>
		</div>
	</Trigger>
</Modal>

<style>
	div {
		background: transparent;
		border: none;
		color: #fff;
		max-width: min-content;
		display: grid;
		grid-template-rows: auto;
		grid-template-columns: auto;
		text-align: center;
		justify-items: center;
		align-items: center;
		margin: 0.25rem;
		cursor: pointer;
	}

	span {
		width: 16rem;
		height: auto;
		aspect-ratio: 2/3;
		margin-bottom: 0.5rem;
	}

	img {
		width: 100%;
		height: 100%;
	}

	h2 {
		text-wrap: auto;
		margin: 0;
		margin-bottom: 0.5rem;
	}
</style>
