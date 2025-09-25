<script lang="ts">
    import { setContext, type Component } from "svelte"
    import views from "$views"
    import Async from "$lib/components/core/Async.svelte"
    import type { View } from "$lib/scripts/core/types.js"
    const components = views as unknown as Record<string, Component>
    let { Name, Props, Render, Align } = $props() as View<Record<string, unknown>>
    const view = $state({ Name, Props, Render, Align })
    setContext("view", view)
    console.log("Name", Name)
</script>

{#each Object.keys(components) as key (key)}
    {#if key === view.Name}
        <Async from={components[key]} properties={view.Props} />
    {/if}
{/each}
