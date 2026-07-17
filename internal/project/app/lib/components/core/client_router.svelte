<script lang="ts">
    import { views } from "$exports.client"
    import { navigate } from "$lib/scripts/core/navigate"
    import { root } from "$lib/scripts/core/root.svelte"
    import type { SvelteComponent } from "svelte"
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
