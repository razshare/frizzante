<script lang="ts">
    import { setContext, type Component } from "svelte"
    import { views } from "$exports.server"
    import type { View } from "$lib/scripts/core/types.js"
    let { Name, Props, Render, Align } = $props() as View<Record<string, unknown>>
    const components = views as unknown as Record<string, Component>
    const view = $state({ Name, Props, Render, Align })
    setContext("view", view)
</script>

{#each Object.keys(components) as key (key)}
    {@const Component = components[key]}
    {#if key === Name}
        <Component {...view.Props} />
    {/if}
{/each}
