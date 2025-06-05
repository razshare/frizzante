<script lang="ts">
    import {setContext} from 'svelte'
    import {views} from "$lib/components/views/client.ts";
    import ClientViewLoader from "$frz/components/ClientViewLoader.svelte";
    import type {View} from "$frz/types.ts";

    let {name, data, error, renderMode} = $props() as View<false>
    const view = $state({name, data, error, renderMode})
    setContext("view", view)
</script>

{#each Object.keys(views) as key(key)}
    {#if key === view.name}
        <ClientViewLoader from={views[key]}/>
    {/if}
{/each}



