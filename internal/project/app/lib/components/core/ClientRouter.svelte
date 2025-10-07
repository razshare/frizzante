<script lang="ts">
    import { setContext, type SvelteComponent } from "svelte"
    import { views } from "$exports.client"
    import Async from "$lib/components/core/Async.svelte"
    import type { View } from "$lib/scripts/core/types.js"
    let { name, props, render, align }: View<Record<string, unknown>> = $props()
    const components = views as unknown as Record<string, ()=>Promise<SvelteComponent>>
    const view: View<Record<string, unknown>> = $state({ name, props, render, align })
    setContext("view", view)
</script>

{#each Object.keys(components) as key (key)}
    {#if key === view.name}
        <Async from={components[key]()} properties={view.props} />
    {/if}
{/each}
