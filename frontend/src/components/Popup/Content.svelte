<script lang="ts">
import { modalID } from "$lib/store";
import { getContext } from "svelte";

const modalId: string = getContext("modalId");
function handleKeydown(e: KeyboardEvent): void {
	if (e.key === "Escape") {
		e.preventDefault();
		$modalID = "";
	}
}
</script>

<svelte:window onkeydown={handleKeydown} />
<div>
	<dialog aria-modal="true" hidden={modalId !== $modalID}>
		<div class="content">
			<button type="button" onclick={() => ($modalID = "")}>
				<span aria-hidden="true">&times;</span>
			</button>
			<slot />
		</div>
	</dialog>
</div>

<style>
:global(#app):has(dialog:not([hidden])) {
	overflow-y: hidden;
}
dialog:not([hidden]) {
	display: flex;
}
dialog {
	position: fixed;
	z-index: 10;
	left: 0;
	top: 0;
	width: 100%;
	height: 100%;
	background: rgba(0, 0, 0, 0.4);
	align-content: center;
	justify-content: center;
	padding: 0;
	border: 0;
	margin: 0;
}
.content {
	position: relative;
	overflow-y: hidden;
	background: transparent;
	background-color: var(--primary);
	padding: 2rem;
	margin: auto;
	height: 80%;
	width: 90%;
	overflow-y: auto;
}
button {
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
dialog {
	color: #fff;
}
</style>
