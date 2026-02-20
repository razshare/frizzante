<script lang="ts">
    import { setContext, type Component } from "svelte"
    import { views } from "$exports.server"
    import type { View } from "$lib/scripts/core/view.d.ts"
    let { name, props, render, align, type }: View<Record<string, unknown>> = $props()
    const components = views as unknown as Record<string, Component>
    // svelte-ignore state_referenced_locally
    const viewStateValue = { name, props, render, align, type, async pin() {} }
    const view: View<Record<string, unknown>> = $state(viewStateValue)
    setContext("view", view)
</script>

{#each Object.keys(components) as key (key)}
    {#if key === view.name}
        {@const Component = components[key]}
        <Component {...view.props} />
    {/if}
{/each}
