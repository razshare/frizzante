<script lang="ts">
    import { setContext } from "svelte"
    import { views } from "$exports/client.ts"
    import ClientViewLoader from "$frizzante/core/components/ClientViewLoader.svelte"
    import type { View } from "$frizzante/core/types.ts"
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    const components = views as Record<string, Component>
    let {
        name,
        props: remoteProps,
        render,
        align,
    } = $props() as View<Record<string, unknown>>
    const view = $state({ name, props: remoteProps, render, align })
    setContext("view", view)
</script>

{#each Object.keys(components) as key (key)}
    {#if key === view.name}
        <ClientViewLoader from={components[key]} properties={view.props} />
    {/if}
{/each}

