<script lang="ts">
    import {setContext} from "svelte";
    import type {View} from "$frz/types.ts";
    import {views} from "$lib/components/views/server.ts";

    let {name, data, error, renderMode} = $props() as View<false>
    const view = $state({name, data, error, renderMode})
    setContext("view", view)
</script>

{#each Object.keys(views) as key(key)}
    {@const Component = views[key]}
    {#if key === name}
        <Component/>
    {/if}
{/each}