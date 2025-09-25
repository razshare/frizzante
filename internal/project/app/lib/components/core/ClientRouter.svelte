<script lang="ts">
    import { setContext } from "svelte"
    import {views as components} from "$exports.client"
    import Async from "$lib/components/core/Async.svelte"
    import type { View } from "$lib/scripts/core/types.js"
    let { Name, Props, Render, Align } = $props() as View<Record<string, unknown>>
    const view = $state({ Name, Props, Render, Align })
    setContext("view", view)
</script>

{#each Object.keys(components) as key (key)}
    {#if key === view.Name}
        <Async from={components[key]} properties={view.Props} />
    {/if}
{/each}
