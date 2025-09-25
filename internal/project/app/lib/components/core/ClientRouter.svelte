<script lang="ts">
    import { setContext, type SvelteComponent } from "svelte"
    import { views } from "$exports.client"
    import Async from "$lib/components/core/Async.svelte"
    import type { View } from "$lib/scripts/core/types.js"
    let { Name, Props, Render, Align } = $props() as View<Record<string, unknown>>
    const components = views as unknown as Record<string, Promise<SvelteComponent>>
    const view = $state({ Name, Props, Render, Align })
    setContext("view", view)
</script>

{#each Object.keys(components) as key (key)}
    {#if key === view.Name}
        <Async from={components[key]} properties={view.Props} />
    {/if}
{/each}
