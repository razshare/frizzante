<script lang="ts">
    import { views } from "$exports.server"
    import { navigate } from "$lib/scripts/core/navigate"
    import { root } from "$lib/scripts/core/root.svelte"
    type Props = {
        name: keyof typeof views
        props: Record<string, unknown>
        type: "" | "default" | "snapshot"
    }
    let { name = $bindable(), props = $bindable(), type = $bindable() }: Props = $props()
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
    {const View = views[root.view.name]}
    {#if View}
        <View {...root.view.ondone()} />
    {/if}
{/if}
