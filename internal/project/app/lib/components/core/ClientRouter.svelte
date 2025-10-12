<script lang="ts">
    import { setContext, type SvelteComponent } from "svelte"
    import { views } from "$exports.client"
    import type { View } from "$lib/scripts/core/types.js"
    let { name, props, render, align }: View<Record<string, unknown>> = $props()
    const components = views as unknown as Record<string, () => Promise<SvelteComponent>>
    const view: View<Record<string, unknown>> = $state({ name, props, render, align, pending: true })
    setContext("view", view)

    let Component: false | SvelteComponent = $state(false)
    let properties: Record<string, unknown> = $state({})

    $effect(function run() {
        if (view.pending) {
            return
        }
        for (const key of Object.keys(components)) {
            if (key == view.name) {
                view.pending = true
                components[key]().then(function run(result) {
                    Component = result
                    properties = view.props
                    view.pending = false
                })
                break
            }
        }
    })
</script>

{#if Component}
    <Component.default {...properties} />
{/if}
