<script lang="ts">
    import { setContext, type SvelteComponent } from "svelte"
    import { views } from "$exports.client"
    import type { View } from "$lib/scripts/core/types.js"
    let { name, props, render, align }: View<Record<string, unknown>> = $props()
    const components = views as unknown as Record<string, () => Promise<SvelteComponent>>
    const view: View<Record<string, unknown>> = $state({
        name,
        props,
        render,
        align,
        pending: false,
        async snapshot() {
            pending.component = await components[view.name]()
            pending.props = $state.snapshot(view.props)
        },
    })
    setContext("view", view)

    const pending = {
        component: false as false | SvelteComponent,
        props: {} as Record<string, unknown>,
    }
</script>

{#each Object.keys(components) as key (key)}
    {#if key === view.name}
        {#await components[key]()}
            {#if pending.component}
                {@const Component = pending.component}
                <Component.default {...pending.props} />
            {/if}
        {:then Component}
            <Component.default {...view.props} />
        {/await}
    {/if}
{/each}
