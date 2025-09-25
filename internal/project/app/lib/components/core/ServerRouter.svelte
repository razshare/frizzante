<script lang="ts">
    import { setContext } from "svelte"
    import {views as components} from "$exports.server"
    import type { View } from "$lib/scripts/core/types.js"
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
