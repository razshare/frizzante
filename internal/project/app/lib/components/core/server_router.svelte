<script lang="ts">
    import { views } from "$exports.server"
    import type { View } from "$lib/scripts/core/view.d.ts"
    import { setContext, type Component } from "svelte"
    let { name, props, render, align, type }: View = $props()
    // svelte-ignore state_referenced_locally
    setContext("view", { name, props, render, align, type })
    // svelte-ignore state_referenced_locally
    const ServerComponent = views[name] as unknown as Component
</script>

<ServerComponent {...props} />
