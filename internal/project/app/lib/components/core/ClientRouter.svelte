<script lang="ts">
    import { setContext, type SvelteComponent } from "svelte"
    import { views } from "$exports.client"
    import type { View } from "$lib/scripts/core/types.js"
    let { name, props, render, align }: View<Record<string, unknown>> = $props()
    const components = views as unknown as Record<string, () => Promise<SvelteComponent>>
    const view: View<Record<string, unknown>> = $state({ name, props, render, align, pending: false })
    setContext("view", view)
</script>

{#each Object.keys(components) as key (key)}
    {#if key === view.name}
        {#await components[key]() then Component}
            <Component.default {...view.props} />
        {/await}
    {/if}
{/each}
