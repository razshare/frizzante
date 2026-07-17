<style>
    .layout {
        position: fixed;
        left: 0;
        right: 0;
        top: 0;
        bottom: 0;
        display: grid;
        justify-content: center;
        align-content: center;
        font-family: "Noto Sans Gothic", serif;
        text-align: center;
        overflow-y: auto;
        overflow-x: hidden;
    }
</style>

<script lang="ts">
    import { root } from "$lib/scripts/core/root.svelte"
    import { type Snippet } from "svelte"
    import { scale } from "svelte/transition"
    type Props = {
        title?: string
        children: Snippet
    }
    let { title = root.view.name, children }: Props = $props()
</script>

<svelte:head>
    <meta charset="UTF-8" />
    <meta
        name="viewport"
        content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0"
    />
    <title>{title}</title>
</svelte:head>

<div in:scale={{ duration: 100 }} class="layout">
    {@render children()}
</div>
