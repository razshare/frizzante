<script lang="ts">
    import { views } from "$exports.client"
    import { navigate } from "$lib/scripts/core/navigate"
    import { root } from "$lib/scripts/core/root.svelte"
    import type { ClientRouterProps } from "$lib/types/core/client_router_props"
    import type { SvelteComponent } from "svelte"
    let { name = $bindable(), props = $bindable(), type = $bindable() }: ClientRouterProps = $props()
    navigate()
    root.type = type
    root.view = {
        name,
        ondone() {
            return props
        },
    }
    let View: false | SvelteComponent = $state(false)
    $effect(function start() {
        if (root.view.name !== "") {
            views[root.view.name]().then(function run(viewLocal) {
                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                //@ts-expect-error
                View = viewLocal
                props = root.view.ondone()
            })
        }
    })
</script>

{#if View}
    <View.default {...props} />
{/if}
