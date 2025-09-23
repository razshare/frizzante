<script lang="ts">
    import { setContext, type Component } from "svelte"
    import type { View } from "$lib/scripts/core/types.js"
    import { views } from "./exports.server.ts"
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    const components = views as Record<string, Component>
    let { Name, Props, Render, Align } = $props() as View<Record<string, unknown>>
    const view = $state({ Name, Props, Render, Align })
    setContext("view", view)
</script>

{#each Object.keys(components) as key (key)}
    {@const Component = components[key]}
    {#if key === Name}
        <Component {...view.Props} />
    {/if}
{/each}
