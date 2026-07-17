<script lang="ts">
    import { views } from "$exports.server"
    import { navigate } from "$lib/scripts/core/navigate"
    import { root } from "$lib/scripts/core/root.svelte"
    import type { ServerRouterProps } from "$lib/types/core/server_router_props"
    import { type Component } from "svelte"
    let { name = $bindable(), props = $bindable(), type = $bindable() }: ServerRouterProps = $props()
    navigate()
    root.type = type
    root.view = {
        name,
        ondone() {
            return props
        },
    }
</script>

{#if root.view.name !== ""}
    {const View = views[root.view.name] as Component<Record<string, unknown>>}
    {#if View}
        <View {...root.view.ondone()} />
    {/if}
{/if}
