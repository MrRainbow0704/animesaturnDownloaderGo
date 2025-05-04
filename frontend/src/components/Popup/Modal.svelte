<script lang="ts">
	import { setContext } from "svelte";

	export let modalId: string;
	setContext("modalId", modalId);
</script>

<slot />
