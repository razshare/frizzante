<script lang="ts">
    import { setContext, type SvelteComponent } from "svelte"
    import { views } from "$exports.client"
    import type { View } from "$lib/scripts/core/view.d.ts"
    let { name, props, render, align, type }: View<Record<string, unknown>> = $props()
    const components = views as unknown as Record<string, () => Promise<SvelteComponent>>
    // svelte-ignore state_referenced_locally
    const viewStateValue = { name, props, render, align, type, pin }
    const view: View<Record<string, unknown>> = $state(viewStateValue)
    setContext("view", view)
    const pending = {
        component: false as false | SvelteComponent,
        props: {} as Record<string, unknown>,
    }
    async function pin() {
        pending.component = await components[view.name]()
        pending.props = $state.snapshot(view.props)
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
