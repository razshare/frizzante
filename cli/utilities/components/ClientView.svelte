<script lang="ts">
    import { setContext } from "svelte"
    import { views } from "$lib/exports/client.ts"
    import ClientViewLoader from "$lib/utilities/components/ClientViewLoader.svelte"
    import type { View } from "$lib/utilities/types.ts"

    let { name, data, error, renderMode } = $props() as View<unknown>
    const view = $state({ name, data, error, renderMode })
    setContext("view", view)
</script>

{#each Object.keys(views) as key (key)}
    {#if key === view.name}
        <ClientViewLoader from={views[key]} />
    {/if}
{/each}
