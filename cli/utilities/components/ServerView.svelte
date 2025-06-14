<script lang="ts">
    import { setContext, type Component } from "svelte"
    import type { View } from "$lib/utilities/types.ts"
    import { views } from "$lib/exports/server.ts"
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    const components = views as Record<string, Component>
    let { name, data, renderMode } = $props() as View<Record<string,unknown>>
    const view = $state({ name, data, renderMode })
    setContext("view", view)
</script>

{#each Object.keys(components) as key (key)}
    {@const Component = components[key]}
    {#if key === name}
        <Component {...view.data}/>
    {/if}
{/each}
