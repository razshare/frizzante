<script lang="ts">
    import { setContext } from "svelte"
    import { views } from "$exports/client.ts"
    import ClientViewLoader from "$frizzante/components/ClientViewLoader.svelte"
    import type { View } from "$frizzante/types.ts"
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    const components = views as Record<string, Component>
    let { name, data, renderMode } = $props() as View<Record<string,unknown>>
    const view = $state({ name, data, renderMode })
    setContext("view", view)
</script>

{#each Object.keys(components) as key (key)}
    {#if key === view.name}
        <ClientViewLoader from={components[key]} properties={view.data}/>
    {/if}
{/each}
