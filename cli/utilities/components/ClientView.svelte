<script lang="ts">
    import { setContext } from "svelte"
    import { views } from "$lib/exports/client.ts"
    import ClientViewLoader from "$lib/utilities/components/ClientViewLoader.svelte"
    import type { View } from "$lib/utilities/types.ts"

    let { name, data, renderMode } = $props() as View<Record<string,unknown>>
    const view = $state({ name, data, renderMode })
    setContext("view", view)
</script>

{#each Object.keys(views) as key (key)}
    {#if key === view.name}
        <ClientViewLoader from={views[key]} properties={view.data}/>
    {/if}
{/each}
