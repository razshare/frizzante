<script lang="ts">
    import { setContext } from "svelte"
    import type { View } from "$lib/utilities/types.ts"
    import { views } from "$lib/exports/server.ts"

    let { name, data, error, renderMode } = $props() as View<unknown>
    const view = $state({ name, data, error, renderMode })
    setContext("view", view)
</script>

{#each Object.keys(views) as key (key)}
    {@const Component = views[key]}
    {#if key === name}
        <Component />
    {/if}
{/each}
