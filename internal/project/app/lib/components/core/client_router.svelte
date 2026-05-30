<script lang="ts">
    import { views } from "$exports.client"
    import type { View } from "$lib/scripts/core/view.d.ts"
    import { setContext, type SvelteComponent } from "svelte"
    let { name, props, render, align, type }: View = $props()
    let ClientComponent = $state(false) as false | SvelteComponent
    // svelte-ignore state_referenced_locally
    let view = $state({ name, props, render, align, type })
    let counter = 0
    setContext("view", view)
    $effect(route)
    function route() {
        const id = ++counter
        const promise = views[view.name]()
        promise
            .then(function resolved(component) {
                // there's a chance the user has clicked a new link
                // and swapped views again while the component was loading [...]
                if (counter != id) {
                    // [...] when that happens we need to bail out,
                    // because most likely the new request asks for a different component,
                    // and we don't want to present the old component view
                    return
                }
                ClientComponent = component as unknown as SvelteComponent
            })
            .catch(function failed(error) {
                console.error(error)
            })
    }
</script>

{#if ClientComponent}
    <ClientComponent.default {...view.props} />
{/if}
