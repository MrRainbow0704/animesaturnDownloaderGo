<script lang="ts">
import { setContext, type Snippet } from "svelte";

let { modalId, children }: { modalId: string; children?: Snippet } = $props();

setContext("modalId", modalId);
</script>

{@render children?.()}
