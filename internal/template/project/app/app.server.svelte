<script lang="ts">
    import { setContext, type Component } from "svelte"
    import type { View } from "$lib/scripts/core/types.ts"
    import { views } from "./exports.server.ts"
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
    {@const Component = components[key]}
    {#if key === name}
        <Component {...view.props} />
    {/if}
{/each}
