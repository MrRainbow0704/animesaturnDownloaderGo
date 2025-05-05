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
		<div style="background-image: url({anime.Info.Poster});">
			<h2>
				{anime.Title}
			</h2>
		</div>
	</Trigger>
</Modal>

<style>
	div {
		background: transparent;
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
</style>
