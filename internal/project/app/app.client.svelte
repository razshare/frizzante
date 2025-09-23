<script lang="ts">
    import { setContext } from "svelte"
    import { views } from "./exports.client.ts"
    import Async from "./app.async.svelte"
    import type { View } from "$lib/scripts/core/types.js"
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    const components = views as Record<string, Component>
    let { Name, Props, Render, Align } = $props() as View<Record<string, unknown>>
    const view = $state({ Name, Props, Render, Align })
    setContext("view", view)
</script>

{#each Object.keys(components) as key (key)}
    {#if key === view.Name}
        <Async from={components[key]} properties={view.Props} />
    {/if}
{/each}
