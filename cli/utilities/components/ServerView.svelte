<script lang="ts">
    import { setContext } from "svelte"
    import type { View } from "$lib/utilities/types.ts"
    import { views } from "$lib/exports/server.ts"

    let { name, data, renderMode } = $props() as View<Record<string,unknown>>
    const view = $state({ name, data, renderMode })
    setContext("view", view)
</script>

{#each Object.keys(views) as key (key)}
    {@const Component = views[key]}
    {#if key === name}
        <Component {...view.data}/>
    {/if}
{/each}
